package bot_misc_cmd

import (
	"github.com/raf924/bot-misc-cmd/internal/pkg"
	"github.com/raf924/bot/pkg/bot"
)

func init() {
	bot.HandleCommand(&pkg.AfkCommand{})
	bot.HandleCommand(&pkg.ListCommand{})
	bot.HandleCommand(&pkg.MathCommand{})
	bot.HandleCommand(&pkg.PingCommand{})
	bot.HandleCommand(&pkg.CmdCommand{})
}
