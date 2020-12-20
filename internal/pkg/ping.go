package pkg

import (
	"fmt"
	"github.com/raf924/bot/api/messages"
	"github.com/raf924/bot/pkg/bot"
	"github.com/raf924/bot/pkg/bot/command"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"time"
)

func init() {
	bot.HandleCommand(&PingCommand{})
}

type PingCommand struct {
	pings map[string]time.Time
	bot   command.Executor
}

func (p *PingCommand) Init(bot command.Executor) error {
	p.bot = bot
	p.pings = map[string]time.Time{}
	return nil
}

func (p *PingCommand) Name() string {
	return "ping"
}

func (p *PingCommand) Aliases() []string {
	return nil
}

func (p *PingCommand) Execute(_ *messages.CommandPacket) ([]*messages.BotPacket, error) {
	id := rand.Int()
	hexId := fmt.Sprintf("%x", id)
	p.pings[hexId] = time.Now()
	return []*messages.BotPacket{
		{
			Timestamp: timestamppb.Now(),
			Message:   hexId,
			Recipient: p.bot.BotUser(),
			Private:   true,
		},
	}, nil
}

func (p *PingCommand) OnChat(message *messages.MessagePacket) ([]*messages.BotPacket, error) {
	if ping, ok := p.pings[message.Message]; ok {
		elapsed := time.Since(ping)
		delete(p.pings, message.Message)
		return []*messages.BotPacket{
			{
				Timestamp: timestamppb.Now(),
				Message:   fmt.Sprintf("Current ping is %s", elapsed.String()),
				Recipient: nil,
				Private:   false,
			},
		}, nil
	}
	return nil, nil
}

func (p *PingCommand) OnUserEvent(_ *messages.UserPacket) ([]*messages.BotPacket, error) {
	return nil, nil
}

func (p *PingCommand) IgnoreSelf() bool {
	return false
}
