package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(genConfigCmd)
}

var genConfigCmd = &cobra.Command{
	Use:   "gen-config",
	Short: "generating example config",
	Long:  `generating example config
usage: redis.exe gen-config > %HOMEPATH%/.config/redis/config.toml`,
	Run: genConfig,
}

func genConfig(cmd *cobra.Command, args []string) {
	fmt.Println("config example blah blah blah")
}

