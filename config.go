package go_redis_server

import (
	"fmt"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server	Server `toml:"server"`
		Redis	Redis `toml:"redis"`
		Log	Log `toml:"log"`
	}

	Server struct {
		Addr 	string `toml:"addr"`
	}

	Redis struct {
		UseSSD bool `toml:"use_ssd" mapstructure:"use_ssd"`
	}

	Log struct {
		LogFolder string `toml:"log_folder" mapstructure:"log_folder"`
	}
)

func LoadConfig(cfgFile string) (Config, error) {
	viper.SetConfigFile(cfgFile)
	err := viper.ReadInConfig()
	fmt.Println("Using config file:", viper.ConfigFileUsed())

	cfg := Config{}
	viper.Unmarshal(&cfg)
	return cfg, err
}
