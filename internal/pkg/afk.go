package pkg

import (
	"fmt"
	"github.com/raf924/bot/pkg/bot/command"
	messages "github.com/raf924/connector-api/pkg/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

type AfkCommand struct {
	afk map[string]string
}

func (a *AfkCommand) Init(_ command.Executor) error {
	a.afk = map[string]string{}
	return nil
}

func (a *AfkCommand) Name() string {
	return "afk"
}

func (a *AfkCommand) Aliases() []string {
	return []string{"away"}
}

func (a *AfkCommand) Execute(command *messages.CommandPacket) ([]*messages.BotPacket, error) {
	a.afk[command.GetUser().GetNick()] = command.GetArgString()
	return []*messages.BotPacket{
		{
			Timestamp: timestamppb.Now(),
			Message:   fmt.Sprintf("User @%s is now AFK", command.GetUser().GetNick()),
			Recipient: nil,
			Private:   false,
		},
	}, nil
}

func (a *AfkCommand) OnChat(message *messages.MessagePacket) ([]*messages.BotPacket, error) {
	var packets []*messages.BotPacket
	_, isAfk := a.afk[message.GetUser().GetNick()]
	if isAfk {
		delete(a.afk, message.GetUser().GetNick())
		packets = append(packets, &messages.BotPacket{
			Timestamp: timestamppb.Now(),
			Message:   "Welcome back",
			Recipient: message.GetUser(),
			Private:   message.GetPrivate(),
		})
	}
	for nick, reason := range a.afk {
		if !strings.Contains(message.GetMessage(), fmt.Sprintf("@%s", nick)) {
			continue
		}
		packets = append(packets, &messages.BotPacket{
			Timestamp: timestamppb.Now(),
			Message:   fmt.Sprintf("%s is afk %s", nick, reason),
			Recipient: message.GetUser(),
			Private:   false,
		})
	}
	return packets, nil
}

func (a *AfkCommand) OnUserEvent(packet *messages.UserPacket) ([]*messages.BotPacket, error) {
	return nil, nil
}

func (a *AfkCommand) IgnoreSelf() bool {
	return true
}
