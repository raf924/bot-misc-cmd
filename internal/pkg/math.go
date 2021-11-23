package pkg

import (
	"github.com/dop251/goja"
	"github.com/raf924/connector-sdk/command"
	"github.com/raf924/connector-sdk/domain"
	"time"
)

var _ command.Command = (*MathCommand)(nil)

const TimeoutMessage = "Exceeded computation timeout"

type MathCommand struct {
	command.NoOpInterceptor
}

func (m *MathCommand) Init(_ command.Executor) error {
	return nil
}

func (m *MathCommand) Name() string {
	return "math"
}

func (m *MathCommand) Aliases() []string {
	return nil
}

func (m *MathCommand) Execute(command *domain.CommandMessage) ([]*domain.ClientMessage, error) {
	r := goja.New()
	valChan := make(chan string)
	go func() {
		val, err := r.RunString(command.ArgString())
		if err != nil {
			val = r.ToValue(err.Error())
		}
		valChan <- val.String()
	}()
	timer := time.NewTimer(5 * time.Second)
	go func() {
		<-timer.C
		r.Interrupt(nil)
		valChan <- TimeoutMessage
	}()
	val := <-valChan
	timer.Stop()
	return []*domain.ClientMessage{
		domain.NewClientMessage(val, command.Sender(), command.Private()),
	}, nil
}
