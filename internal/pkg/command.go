package pkg

import (
	"fmt"
	"github.com/raf924/bot/pkg/bot/command"
	"github.com/raf924/bot/pkg/bot/permissions"
	"github.com/raf924/bot/pkg/domain"
	"github.com/raf924/bot/pkg/storage"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var _ command.Command = (*CmdCommand)(nil)

type setter string

const (
	setText setter = "set-txt"
	setJs   setter = "set-js"
)

var commandNameRegex = regexp.MustCompile("(?i)([a-zA-Z][a-zA-Z0-9-_]*)")
var argRegex = regexp.MustCompile("(?i)(%[0-9]+%)")

type Creator struct {
	Nick string `json:"nick"`
	Id   string `json:"id"`
}

func (c *Creator) Is(user *domain.User) bool {
	return c.Id == "" && user.Nick() == c.Nick || user.Id() == c.Id
}

type customCommand struct {
	Setter  setter   `json:"setter"`
	Creator *Creator `json:"creator"`
	Source  string   `json:"source"`
}

type CmdCommand struct {
	command.NoOpInterceptor
	storage  storage.Storage
	commands map[string]customCommand
	bot      command.Executor
}

func (c *CmdCommand) Init(bot command.Executor) error {
	log.Println("Init cmd Command")
	commandsLocation := bot.ApiKeys()["commandsLocation"]
	cmdStorage, err := storage.NewFileStorage(commandsLocation)
	if err != nil {
		log.Println(err)
		cmdStorage = storage.NewNoOpStorage()
	}
	c.storage = cmdStorage
	c.bot = bot
	c.commands = map[string]customCommand{}
	c.load()
	return nil
}

func (c *CmdCommand) Name() string {
	return "cmd"
}

func (c *CmdCommand) Aliases() []string {
	return nil
}

func runCommand(cmd *customCommand, botCommand *domain.CommandMessage) (string, error) {
	switch cmd.Setter {
	case setText:
		return runTextCommand(cmd.Source, botCommand), nil
	}
	return "", fmt.Errorf("")
}

func runTextCommand(source string, botCommand *domain.CommandMessage) string {
	args := botCommand.Args()
	commandContent := strings.ReplaceAll(source, "%sender%", botCommand.Sender().Nick())
	return argRegex.ReplaceAllStringFunc(commandContent, func(s string) string {
		argIndexStr := argRegex.FindAllStringSubmatch(s, -1)[0][1]
		argIndex, err := strconv.ParseInt(argIndexStr, 10, 64)
		if err != nil {
			return s
		}
		return args[argIndex]
	})
}

func (c *CmdCommand) Execute(command *domain.CommandMessage) ([]*domain.ClientMessage, error) {
	set := command.Args()[0]
	cmdName := command.Args()[1]
	if !commandNameRegex.MatchString(cmdName) {
		return nil, fmt.Errorf("command name doesn't match approved pattern")
	}
	defer c.save()
	if set == "unset" {
		var message string
		if c.unsetCommand(cmdName, command.Sender()) {
			message = "unset command %s"
		} else {
			message = "couldn't unset command %s"
		}
		return []*domain.ClientMessage{
			domain.NewClientMessage(fmt.Sprintf(message, cmdName), command.Sender(), command.Private()),
		}, nil
	}
	switch setter(set) {
	case setText:
		fallthrough
	case setJs:
		couldSet := c.setCommand(cmdName, customCommand{
			Setter: setter(set),
			Creator: &Creator{
				Nick: command.Sender().Nick(),
				Id:   command.Sender().Id(),
			},
			Source: strings.Split(command.ArgString(), fmt.Sprintf("%s %s", set, cmdName))[1],
		})
		var message string
		if couldSet {
			message = "set command %s"
		} else {
			message = "could not set command %s"
		}
		return []*domain.ClientMessage{
			domain.NewClientMessage(fmt.Sprintf(message, cmdName), command.Sender(), command.Private()),
		}, nil
	}
	return nil, nil
}

func (c *CmdCommand) OnChat(message *domain.ChatMessage) ([]*domain.ClientMessage, error) {
	log.Println("OnChat", message.Message())
	if len(c.bot.Trigger()) == 0 {
		return nil, nil
	}
	if !strings.HasPrefix(message.Message(), c.bot.Trigger()) {
		return nil, nil
	}
	argString := strings.TrimPrefix(message.Message(), c.bot.Trigger())
	args := strings.Split(argString, " ")
	if len(args) == 0 || len(args[0]) == 0 {
		return nil, nil
	}
	possibleCommand := args[0]
	cmd, ok := c.commands[possibleCommand]
	if !ok {
		return nil, nil
	}
	argString = strings.TrimSpace(strings.TrimPrefix(argString, possibleCommand))
	commandPacket := domain.NewCommandMessage(possibleCommand, args[1:], argString, message.Sender(), message.Private(), message.Timestamp())
	text, err := runCommand(&cmd, commandPacket)
	if err != nil {
		return nil, err
	}
	return []*domain.ClientMessage{
		domain.NewClientMessage(text, message.Sender(), message.Private()),
	}, nil
}

func (c *CmdCommand) unsetCommand(cmdName string, unsetter *domain.User) bool {
	var cmd customCommand
	var ok bool
	if cmd, ok = c.commands[cmdName]; !ok {
		return false
	}
	if !c.bot.UserHasPermission(unsetter, permissions.MOD) {
		if cmd.Creator == nil || cmd.Creator.Is(unsetter) {
			return false
		}
	}
	delete(c.commands, cmdName)
	return true
}

func (c *CmdCommand) setCommand(cmdName string, cmd customCommand) bool {
	if _, ok := c.commands[cmdName]; ok {
		return false
	}
	c.commands[cmdName] = cmd
	return true
}

func (c *CmdCommand) load() {
	err := c.storage.Load(&c.commands)
	if err != nil {
		log.Println("error reading commands file:", err.Error())
		return
	}
}

func (c *CmdCommand) save() {
	c.storage.Save(c.commands)
}
