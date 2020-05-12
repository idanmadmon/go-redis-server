package go_redis_server

import (
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
)

var cfg Redis

func initialize(config Config) {
	cfg = config.Redis
	DB = safeDB{make(map[string]string, 0), &sync.Mutex{}}
	setCommands()
}

func Start(cfg Config) error {
	initialize(cfg)
	log.WithField("config", cfg).Info("Application Initialize")
	defer stop()
	c, err := ParseRequest("*3\r\n$3\r\nset\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")
	if err != nil {
		log.Debug(err)
	} else {
		log.Debug(strconv.Itoa(len(c)) + "[" + *c[0] + "," + *c[1] + "," + *c[2] + "]")
	}
	return nil
	//return listen(cfg.Server.Addr, cfg.Server.ConnType)
}

func stop() {
	log.Info("Start stopping sequence")

	//wait for final actions to finish
	DB.Lock()
	DB.Unlock()
}

func listen(host, connType string) error {
	l, err := net.Listen(connType, host)
	if err != nil {
		log.WithError(err).Fatal("Error listening")
		return err
	}
	defer l.Close()

	log.WithField("host", host).WithField("type", connType).Info("Start listening")
	interrupted := false
	go accept(l, &interrupted)
	waitForInterrupt(l, &interrupted)
	return nil
}

func waitForInterrupt(l net.Listener, interrupted *bool) {
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, os.Interrupt)
	<-signalC
	*interrupted = true
	l.Close()
}

func accept(l net.Listener, interrupted *bool) {
	for {

		conn, err := l.Accept()
		if err != nil {
			if !*interrupted {
				log.WithError(err).Error("Error accepting")
			}
			return
		}

		go handleRequest(conn)
	}
}