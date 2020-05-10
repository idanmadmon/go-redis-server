package cmd

import (
	redis "github.com/idan/go-redis-server"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	logLevel     string

	rootCmd = &cobra.Command{
		Use:   "go-redis-server",
		Short: "A simple redis server",
		Long: `A simple redis server written in go, has only get and set functions.`,
		Run: func(cmd *cobra.Command, args []string) {
			redis.Run()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "Info", "log level")
}

func initConfig() {
}