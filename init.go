package bot_misc_cmd

import (
	"github.com/raf924/bot-misc-cmd/v2/internal/pkg"
	"github.com/raf924/connector-sdk/command"
)

func init() {
	command.HandleCommand(&pkg.AfkCommand{})
	command.HandleCommand(&pkg.ListCommand{})
	command.HandleCommand(&pkg.MathCommand{})
	command.HandleCommand(&pkg.PingCommand{})
	command.HandleCommand(&pkg.CmdCommand{})
}
