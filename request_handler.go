package go_redis_server

import (
	"errors"
	"net"

	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var requestc chan Request

type (
	Request struct {
		net.Conn
	}

	RequestHandle struct {
		Clients *map[uuid.UUID]net.Conn
		Worker
	}
)

func (r *RequestHandle) Start() error {
	if r.Clients == nil {
		return	errors.New("request handler can't start without given clients")
	}

	if requestc == nil {
		requestc = make(chan Request, 0)
	}

	for i := 0; i < cfg.RequestWorkers; i++ {
		go r.run()
	}
	return nil
}

func (r *RequestHandle) run() {
	for !*r.Interrupt {
		conn := <-requestc
		msg, id, err := HandleRequest(conn)
		(*r.Clients)[id] = conn
		if err != nil {
			ReplyError(err, id)
			continue
		}
		parsec <- msg
	}
}

func HandleRequest(conn net.Conn) (Message, uuid.UUID, error){
	id := uuid.Must(uuid.NewV4())

	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	if err != nil {
		log.WithError(err).Error("Error reading request")
		return Message{}, id, err
	}

	log.WithField("size", reqLen).Info("Got new request")
	return Message{id, string(buf)}, id, nil
}