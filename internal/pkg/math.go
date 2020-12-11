package pkg

import (
	"github.com/dop251/goja"
	"github.com/raf924/bot/api/messages"
	"github.com/raf924/bot/pkg/bot"
	"github.com/raf924/bot/pkg/bot/command"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func init() {
	bot.HandleCommand(&MathCommand{})
}

type MathCommand struct {
	command.NoOpInterceptor
}

func (m *MathCommand) Init(bot command.Executor) error {
	return nil
}

func (m *MathCommand) Name() string {
	return "math"
}

func (m *MathCommand) Aliases() []string {
	return nil
}

func (m *MathCommand) Execute(command *messages.CommandPacket) ([]*messages.BotPacket, error) {
	r := goja.New()
	val, err := r.RunString(command.GetArgString())
	if err != nil {
		val = r.ToValue(err.Error())
	}
	return []*messages.BotPacket{
		{
			Timestamp: timestamppb.Now(),
			Message:   val.String(),
			Recipient: command.GetUser(),
			Private:   command.GetPrivate(),
		},
	}, nil
}
