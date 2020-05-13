package go_redis_server

import (
	"fmt"
	"io"
	"os"

	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/t-tomalak/logrus-prefixed-formatter"
)

type (
	Config struct {
		Server	Server	`toml:"server"`
		Redis	Redis 	`toml:"redis"`
		Log		Log 	`toml:"log"`
	}

	Server struct {
		Addr 		string `toml:"addr"`
		ConnType	string `toml:"conn_type" mapstructure:"conn_type"`
	}

	Redis struct {
		UseSSD 				bool	`toml:"use_ssd" mapstructure:"use_ssd"`
		DisableOverride 	bool	`toml:"disable_override" mapstructure:"disable_override"`
		RequestWorkers		int 	`toml:"request_workers" mapstructure:"request_workers"`
		CommandsWorkers		int 	`toml:"commands_workers" mapstructure:"commands_workers"`
		ParseWorkers		int 	`toml:"parse_workers" mapstructure:"parse_workers"`
		ReplyWorkers		int 	`toml:"reply_workers" mapstructure:"reply_workers"`
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

func InitLogger(dir, level string) error {
	lvl, err := log.ParseLevel(level)
	if err != nil {
		return err
	}
	log.SetLevel(lvl)

	log.SetFormatter(&prefixed.TextFormatter{
		DisableColors: true,
		TimestampFormat : "2006-01-02 15:04:05.000000",
		FullTimestamp:true,
		ForceFormatting: true,
	})

	ljack := &lumberjack.Logger{
		Filename:   "logs/go-redis.log",
		MaxSize:    20,		// megabytes
		MaxBackups: 50,
		MaxAge:     30,		//days
		Compress:   false,	// disabled by default
	}
	mw := io.MultiWriter(os.Stdout, ljack)
	log.SetOutput(mw)
	return nil
}
