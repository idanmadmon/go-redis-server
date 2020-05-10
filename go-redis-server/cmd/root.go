package cmd

import (
	redis "github.com/idan/go-redis-server"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"strings"
)

var (
	// Used for flags.
	logLevel	string
	cfgFile		string
	cfg			redis.Config

	rootCmd = &cobra.Command{
		Use:   "go-redis-server",
		Short: "A simple redis server",
		Long: `A simple redis server written in go, has only get and set functions.`,
		PreRun: initConfig,
		Run: func(cmd *cobra.Command, args []string) {
			redis.Run(cfg)
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "Info", "log level")
	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file Default(%HOMEPATH%/.config/redis/.)")
}

func initConfig(cmd *cobra.Command, args []string) {
	if cfgFile == "" || strings.LastIndex(cfgFile, ".toml") != len(cfgFile) - 5 {
		home, err := homedir.Dir()
		exitOnErr(err)
		cfgFile = home + "/.config/redis/config.toml"
	}

	var err error = nil
	cfg, err = redis.LoadConfig(cfgFile)
	exitOnErr(err)
}