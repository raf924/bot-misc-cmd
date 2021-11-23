package pkg

import (
	"github.com/raf924/connector-sdk/command"
	"github.com/raf924/connector-sdk/domain"
	"sort"
	"strings"
)

var _ command.Command = (*ListCommand)(nil)

type ListCommand struct {
	bot command.Executor
	command.NoOpInterceptor
}

func (l *ListCommand) Init(bot command.Executor) error {
	l.bot = bot
	return nil
}

func (l *ListCommand) Name() string {
	return "list"
}

func (l *ListCommand) Aliases() []string {
	return []string{"l"}
}

func (l *ListCommand) Execute(command *domain.CommandMessage) ([]*domain.ClientMessage, error) {
	onlineUsers := l.bot.OnlineUsers().All()
	users := make([]string, 0, len(onlineUsers))
	sort.Sort(domain.Users(onlineUsers))
	for _, u := range onlineUsers {
		users = append(users, u.Nick())
	}
	return []*domain.ClientMessage{
		domain.NewClientMessage(strings.Join(users, ", "), command.Sender(), command.Private()),
	}, nil
}
