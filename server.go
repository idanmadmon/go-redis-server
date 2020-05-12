package go_redis_server

import (
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"sync"
)

var cfg Redis

func initialize(config Config) {
	cfg = config.Redis
	DB = safeDB{make(map[string]string, 0), &sync.Mutex{}}
}

func Start(cfg Config) error {
	initialize(cfg)
	log.WithField("config", cfg).Info("Application Initialize")
	interrupted := false
	h := Handler{cfg.Redis,&interrupted}
	p := Parse{cfg.Redis,&interrupted}
	p.Start()
	c := Commands{nil, cfg.Redis,&interrupted}
	c.Start()
	defer stop()
	return listen(cfg.Server.Addr, cfg.Server.ConnType, h, &interrupted)
}

func stop() {
	log.Info("Start stopping sequence")

	//wait for final actions to finish
	DB.Lock()
	DB.Unlock()
}

func listen(host, connType string, h Handler, interrupted *bool) error {
	l, err := net.Listen(connType, host)
	if err != nil {
		log.WithError(err).Fatal("Error listening")
		return err
	}
	defer l.Close()

	log.WithField("host", host).WithField("type", connType).Info("Start listening")
	go accept(l, h, interrupted)
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

func accept(l net.Listener, h Handler, interrupted *bool) {
	for {
		conn, err := l.Accept()
		if err != nil {
			if !*interrupted {
				log.WithError(err).Error("Error accepting")
			}
			return
		}

		go h.HandleRequest(conn)
	}
}