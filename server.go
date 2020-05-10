package go_redis_server

import "fmt"

func Run(cfg Config) {
	fmt.Println("yoooooo")
	fmt.Println(cfg.Log.LogFolder)
	fmt.Println(cfg.Server.Addr)
	fmt.Println(cfg.Redis.UseSSD)
}
