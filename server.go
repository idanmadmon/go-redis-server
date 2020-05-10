package go_redis_server

import (
	log "github.com/sirupsen/logrus"
)

func Run(cfg Config) {
	log.WithField("config", cfg).Info("Application Initialize")
}
