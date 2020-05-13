package go_redis_server

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var cmdsc chan Command

type (
	Commands struct {
		Cmds	map[string]func(args []*string) (string, error)
		db		*DB
		Worker
	}

	Command struct {
		id		uuid.UUID
		data	[]*string
	}
)

func (c *Commands) initialize() {
	cmdsc = make(chan Command, 0)
	c.Cmds = make(map[string]func(args []*string) (string, error), 0)
	c.Cmds["ping"] = pingCommand
	c.Cmds["set"] = setCommand
	c.Cmds["setnx"] = setCommand //TODO
	c.Cmds["get"] = getCommand
}

func (c *Commands) Start() error {
	if cmdsc == nil || c.Cmds == nil {
		c.initialize()
	}

	for i := 0; i < cfg.CommandsWorkers; i++ {
		go c.run()
	}
	return nil
}

func (c *Commands) run() {
	for !*c.Interrupt {
		cmd := <-cmdsc
		f, ok := c.Cmds[*cmd.data[0]]
		if !ok {
			log.WithField("command", *cmd.data[0]).Error("unknown command")
			ReplyError(errors.New("unknown command '" + *cmd.data[0] + "'"), cmd.id)
			continue
		}

		r, _ := f(cmd.data[1:])
		ReplyMessage(r, cmd.id)
	}
}

func pingCommand(args []*string) (string, error) {
	return "+PONG\r\n", nil
}

func setCommand(args []*string) (string, error) {
	return "+OK\r\n", nil
}

func getCommand(args []*string) (string, error) {
	return "$4\r\n1234\r\n", nil
}