package go_redis_server

import log "github.com/sirupsen/logrus"

func exitOnError(err error) {
	log.WithField("Error", err).Fatal("Exiting...")
	//TODO: start workers stop sequence
	panic(err)
}
