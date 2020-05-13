package go_redis_server

import (
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"sync"
)

var cfg Redis

type Worker struct {
	Cfg			Redis
	Interrupt	*bool
}

func initialize(config Config) {
	cfg = config.Redis
	DB = safeDB{make(map[string]string, 0), &sync.Mutex{}}
}

func Start(cfg Config) error {
	initialize(cfg)
	log.WithField("config", cfg).Info("Application Initialize")
	interrupted := false
	//maybe use
	w := Worker{cfg.Redis,&interrupted}
	clients := make(map[uuid.UUID]net.Conn, 0)
	h := RequestHandle{&clients, w}
	h.Start()
	p := Parse{w}
	p.Start()
	c := Commands{nil, w}
	c.Start()
	r := ReplyHandle{&clients, w}
	r.Start()
	defer stop()
	return listen(cfg.Server.Addr, cfg.Server.ConnType, &interrupted)
}

func stop() {
	log.Info("Start stopping sequence")

	//wait for final actions to finish
	DB.Lock()
	DB.Unlock()
}

func listen(host, connType string, interrupted *bool) error {
	l, err := net.Listen(connType, host)
	if err != nil {
		log.WithError(err).Fatal("Error listening")
		return err
	}
	defer l.Close()

	log.WithField("host", host).WithField("type", connType).Info("Start listening")
	go accept(l, interrupted)
	waitForInterrupt(l, interrupted)
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

		requestc <- Request{conn}
	}
}