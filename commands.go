package go_redis_server

import (
	uuid "github.com/satori/go.uuid"
)

var cmdsc chan Command

type (
	Commands struct {
		Cmds		map[string]func(args []*string) (string, error)
		Cfg			Redis
		Interrupt	*bool
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
		r, _ := c.Cmds[*cmd.data[0]](cmd.data[1:])
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