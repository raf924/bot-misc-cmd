package pkg

import (
	"fmt"
	"github.com/raf924/connector-sdk/command"
	"github.com/raf924/connector-sdk/domain"
	"strings"
)

var _ command.Command = (*AfkCommand)(nil)

type AfkCommand struct {
	command.NoOpInterceptor
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

func (a *AfkCommand) Execute(command *domain.CommandMessage) ([]*domain.ClientMessage, error) {
	a.afk[command.Sender().Nick()] = command.ArgString()
	return []*domain.ClientMessage{
		domain.NewClientMessage(fmt.Sprintf("User @%s is now AFK", command.Sender().Nick()), nil, false),
	}, nil
}

func (a *AfkCommand) OnChat(message *domain.ChatMessage) ([]*domain.ClientMessage, error) {
	var packets []*domain.ClientMessage
	_, isAfk := a.afk[message.Sender().Nick()]
	if isAfk {
		delete(a.afk, message.Sender().Nick())
		packets = append(packets, domain.NewClientMessage("Welcome back", message.Sender(), message.Private()))
	}
	for nick, reason := range a.afk {
		if !strings.Contains(message.Message(), fmt.Sprintf("@%s", nick)) {
			continue
		}
		packets = append(packets, domain.NewClientMessage(fmt.Sprintf("%s is afk %s", nick, reason), message.Sender(), false))
	}
	return packets, nil
}

func (a *AfkCommand) IgnoreSelf() bool {
	return true
}
