package go_redis_server

var commands map[string]func(args []string)

func setCommands() {
	commands = make(map[string]func(args []string), 0)
}