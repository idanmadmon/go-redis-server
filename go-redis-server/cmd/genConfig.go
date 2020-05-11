package cmd

import (
	"fmt"
	"github.com/pelletier/go-toml"

	redis "github.com/idan/go-redis-server"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(genConfigCmd)
}

var genConfigCmd = &cobra.Command{
	Use:   "gen-config",
	Short: "generating example config",
	Long:  `generating example config
usage: redis.exe gen-config > %HOMEPATH%/.config/go-redis/config.toml`,
	Run: genConfig,
}

func genConfig(cmd *cobra.Command, args []string) {
	cfg := redis.Config{Server: redis.Server{Addr: "localhost:6379", ConnType: "tcp"}, Redis: redis.Redis{}, Log: redis.Log{LogFolder: "logs"}}
	cfgb, err := toml.Marshal(cfg)
	exitOnErr(err)
	fmt.Println(string(cfgb))
}

