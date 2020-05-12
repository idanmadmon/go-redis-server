package go_redis_server

import (
	"reflect"
)

var commands map[string]func(args []string) (string, error)

type Command struct{}

func setCommands() {
	commands = make(map[string]func(args []string) (string, error), 0)
	//TODO: think if this is the best way (get the commands from outside) or not exposing it
	// and implement it like: commands["ping"] = PingCommand
	commands["ping"] = reflect.ValueOf(Command{}).MethodByName(cfg.Ping).Interface().(func(args []string) (string, error))
	commands["set"] = reflect.ValueOf(Command{}).MethodByName(cfg.Set).Interface().(func(args []string) (string, error))
	commands["get"] = reflect.ValueOf(Command{}).MethodByName(cfg.Get).Interface().(func(args []string) (string, error))
}

func (Command) PingCommand(args []string) (string, error) {
	return "PONG", nil
}

func (Command) SetCommand(args []string) (string, error) {
	return "OK", nil
}

func (Command) GetCommand(args []string) (string, error) {
	return "1234", nil
}