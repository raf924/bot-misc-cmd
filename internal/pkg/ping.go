package pkg

import (
	"fmt"
	"github.com/raf924/connector-sdk/command"
	"github.com/raf924/connector-sdk/domain"
	"math/rand"
	"time"
)

var _ command.Command = (*PingCommand)(nil)

type PingCommand struct {
	command.NoOpInterceptor
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

func (p *PingCommand) Execute(_ *domain.CommandMessage) ([]*domain.ClientMessage, error) {
	rand.Seed(time.Now().UnixNano())
	id := rand.Int()
	hexId := fmt.Sprintf("%x", id)
	p.pings[hexId] = time.Now()
	return []*domain.ClientMessage{
		domain.NewClientMessage(hexId, p.bot.BotUser(), true),
	}, nil
}

func (p *PingCommand) OnChat(message *domain.ChatMessage) ([]*domain.ClientMessage, error) {
	if ping, ok := p.pings[message.Message()]; ok {
		elapsed := time.Since(ping)
		delete(p.pings, message.Message())
		return []*domain.ClientMessage{
			domain.NewClientMessage(fmt.Sprintf("Current ping is %s", elapsed.String()), nil, false),
		}, nil
	}
	return nil, nil
}

func (p *PingCommand) IgnoreSelf() bool {
	return false
}
