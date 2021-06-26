package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/raf924/bot/pkg/bot/command"
	"github.com/raf924/bot/pkg/bot/permissions"
	messages "github.com/raf924/connector-api/pkg/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type setter string

const (
	setText = "set-txt"
	setJs   = "set-js"
)

var commandNameRegex = regexp.MustCompile("(?i)([a-zA-Z][a-zA-Z0-9-_]*)")
var argRegex = regexp.MustCompile("(?i)(%[0-9]+%)")

type customCommand struct {
	Setter  setter         `json:"setter"`
	Creator *messages.User `json:"creator"`
	Source  string         `json:"source"`
}

type CmdCommand struct {
	command.NoOpInterceptor
	commandsLocation string
	storageMutex     *sync.Mutex
	commands         map[string]customCommand
	bot              command.Executor
}

func (c *CmdCommand) Init(bot command.Executor) error {
	c.commandsLocation = bot.ApiKeys()["commandsLocation"]
	c.storageMutex = &sync.Mutex{}
	c.bot = bot
	c.load()
	return nil
}

func (c *CmdCommand) Name() string {
	return "cmd"
}

func (c *CmdCommand) Aliases() []string {
	return nil
}

func runCommand(cmd *customCommand, botCommand *messages.CommandPacket) (string, error) {
	switch cmd.Setter {
	case setText:
		return runTextCommand(cmd.Source, botCommand), nil
	}
	return "", fmt.Errorf("")
}

func runTextCommand(source string, botCommand *messages.CommandPacket) string {
	args := botCommand.GetArgs()[1:]
	commandContent := strings.ReplaceAll(source, "%sender%", botCommand.GetUser().GetNick())
	return argRegex.ReplaceAllStringFunc(commandContent, func(s string) string {
		argIndexStr := argRegex.FindAllStringSubmatch(s, -1)[0][1]
		argIndex, err := strconv.ParseInt(argIndexStr, 10, 64)
		if err != nil {
			return s
		}
		return args[argIndex]
	})
}

func (c *CmdCommand) Execute(command *messages.CommandPacket) ([]*messages.BotPacket, error) {
	set := command.GetArgs()[0]
	cmdName := command.GetArgs()[1]
	if !commandNameRegex.MatchString(cmdName) {
		return nil, fmt.Errorf("command name doesn't match approved pattern")
	}
	defer c.save()
	if set == "unset" {
		var message string
		if c.unsetCommand(cmdName, command.GetUser()) {
			message = "unset command %s"
		} else {
			message = "couldn't unset command %s"
		}
		return []*messages.BotPacket{
			{
				Timestamp: timestamppb.Now(),
				Message:   fmt.Sprintf(message, cmdName),
				Recipient: command.GetUser(),
				Private:   command.GetPrivate(),
			},
		}, nil
	}
	switch set {
	case setText:
		fallthrough
	case setJs:
		couldSet := c.setCommand(cmdName, customCommand{
			Setter:  setter(set),
			Creator: command.GetUser(),
			Source:  strings.Split(command.GetArgString(), fmt.Sprintf("%s %s", set, cmdName))[1],
		})
		var message string
		if couldSet {
			message = "set command %s"
		} else {
			message = "could not set command %s"
		}
		return []*messages.BotPacket{
			{
				Timestamp: timestamppb.Now(),
				Message:   fmt.Sprintf(message, cmdName),
				Recipient: command.GetUser(),
				Private:   command.GetPrivate(),
			},
		}, nil
	}
	return nil, nil
}

func (c *CmdCommand) OnChat(message *messages.MessagePacket) ([]*messages.BotPacket, error) {
	if len(c.bot.Trigger()) == 0 {
		return nil, nil
	}
	if !strings.HasPrefix(message.GetMessage(), c.bot.Trigger()) {
		return nil, nil
	}
	argString := strings.TrimPrefix(message.GetMessage(), c.bot.Trigger())
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
	commandPacket := &messages.CommandPacket{
		Timestamp: message.GetTimestamp(),
		Command:   possibleCommand,
		Args:      args[1:],
		ArgString: argString,
		User:      message.GetUser(),
		Private:   message.GetPrivate(),
	}
	text, err := runCommand(&cmd, commandPacket)
	if err != nil {
		return nil, err
	}
	return []*messages.BotPacket{
		{
			Timestamp: timestamppb.Now(),
			Message:   text,
			Recipient: message.GetUser(),
			Private:   message.GetPrivate(),
		},
	}, nil
}

func (c *CmdCommand) unsetCommand(cmdName string, unsetter *messages.User) bool {
	var cmd customCommand
	var ok bool
	if cmd, ok = c.commands[cmdName]; !ok {
		return false
	}
	if !c.bot.UserHasPermission(unsetter, permissions.MOD) {
		if cmd.Creator == nil || unsetter.String() != cmd.Creator.String() {
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
	file, err := os.Open(c.commandsLocation)
	if err != nil {
		log.Println("error opening message file:", err.Error())
		return
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&c.commands)
	if err != nil {
		log.Println("error reading message file:", err.Error())
		return
	}
}

func (c *CmdCommand) save() {
	go func() {
		c.storageMutex.Lock()
		defer c.storageMutex.Unlock()
		file, err := os.OpenFile(c.commandsLocation, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Println("error opening command file:", err.Error())
			return
		}
		defer file.Close()
		err = json.NewEncoder(file).Encode(c.commands)
		if err != nil {
			log.Println("error writing to command file:", err.Error())
			return
		}
	}()
}
