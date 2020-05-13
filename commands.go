package go_redis_server

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"strconv"
)

var cmdsc chan Command

type (
	Commands struct {
		Cmds	map[string]func(args []string) (string, error)
		db		*DB
		Worker
	}

	Command struct {
		id		uuid.UUID
		data	[]string
	}
)

func (c *Commands) initialize() {
	cmdsc = make(chan Command, 0)
	c.Cmds = make(map[string]func(args []string) (string, error), 0)
	c.Cmds["ping"] = c.pingCommand
	c.Cmds["set"] = c.setCommand
	c.Cmds["setnx"] = c.setCommand //TODO
	c.Cmds["get"] = c.getCommand
}

func (c *Commands) Start() error {
	if cmdsc == nil || c.Cmds == nil {
		c.initialize()
	}

	for i := 0; i < c.Cfg.CommandsWorkers; i++ {
		go c.run()
	}
	return nil
}

func (c *Commands) run() {
	for !*c.Interrupt {
		cmd := <-cmdsc
		f, ok := c.Cmds[cmd.data[0]]
		if !ok {
			log.WithField("command", cmd.data[0]).Error("unknown command")
			ReplyError(errors.New("unknown command '" + cmd.data[0] + "'"), cmd.id)
			continue
		}

		r, err := f(cmd.data[1:])
		if err != nil {
			ReplyError(err, cmd.id)
			continue
		}

		ReplyMessage(r, cmd.id)
	}
}

func (c *Commands) pingCommand(args []string) (string, error) {
	return "+PONG\r\n", nil
}

func (c *Commands) setCommand(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("bad index")
	}

	err := c.db.set(args[0], args[1])
	if err != nil {
		return "", err
	}

	return "+OK\r\n", nil
}

func (c *Commands) getCommand(args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("bad index")
	}

	val, err := c.db.get(args[0])
	if err != nil {
		//key not found
		return buildRespNullBulkString(), nil
	}

	result := "$" + strconv.Itoa(len(val)) + "\r\n" + val + "\r\n"
	return result, nil
}