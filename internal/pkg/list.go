package pkg

import (
	"github.com/raf924/bot/pkg/bot/command"
	messages "github.com/raf924/connector-api/pkg/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

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

func (l *ListCommand) Execute(command *messages.CommandPacket) ([]*messages.BotPacket, error) {
	var users []string
	for _, u := range l.bot.OnlineUsers() {
		users = append(users, u.Nick)
	}
	return []*messages.BotPacket{
		{
			Timestamp: timestamppb.Now(),
			Message:   strings.Join(users, ", "),
			Recipient: command.GetUser(),
			Private:   command.GetPrivate(),
		},
	}, nil
}
