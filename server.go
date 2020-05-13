package go_redis_server

import (
	"net"
	"os"
	"os/signal"
	"sync"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)


type (
	Worker struct {
		Cfg			Redis
		Interrupt	*bool
	}

	Server struct {
		cfg	Redis
		db	*DB

		h	RequestHandle
		p	Parse
		c	Commands
		r	ReplyHandle
	}
)

func (s *Server) initialize(config Config, interrupted *bool) {
	s.cfg = config.Redis
	w := Worker{s.cfg,interrupted}
	clients := make(map[uuid.UUID]net.Conn, 0)
	db := DB{safeMap{make(map[string]string, 0), &sync.Mutex{}}, s.cfg}
	s.db = &db
	s.h = RequestHandle{&clients, w}
	s.p = Parse{w}
	s.c = Commands{nil, &db, w}
	s.r = ReplyHandle{&clients, w}
}

func (s *Server) Start(cfg Config) error {
	interrupted := false
	s.initialize(cfg, &interrupted)
	log.WithField("config", cfg).Info("Application Initialize")
	s.h.Start()
	s.p.Start()
	s.c.Start()
	s.r.Start()
	defer s.Stop()
	return listen(cfg.Host.Addr, cfg.Host.ConnType, &interrupted)
}

func (s *Server) Stop() {
	log.Info("Start stopping sequence")

	//wait for final actions to finish
	s.db.Lock()
	s.db.Unlock()
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