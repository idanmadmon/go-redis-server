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

	for i := 0; i < r.Cfg.RequestWorkers; i++ {
		go r.run()
	}
	return nil
}

func (r *RequestHandle) run() {
	for !*r.Interrupt {
		conn := <-requestc
		id := uuid.Must(uuid.NewV4())
		(*r.Clients)[id] = conn
		go r.HandleRequest(id, conn)
	}
}

func (r *RequestHandle) HandleRequest(id uuid.UUID, conn net.Conn) {
	var err error = nil
	msg := ""
	for err == nil{
		msg, err = readRequest(conn)
		if err != nil {
			ReplyError(err, id)
			continue
		}
		parsec <- Message{id, msg}
	}

	conn.Close()
	delete(*r.Clients,id)
}

func readRequest(conn net.Conn) (string, error){
	r := ""
	bytesToRead := 1024
	reqLen := bytesToRead
	var err error = nil

	for reqLen == bytesToRead {
		buf := make([]byte, bytesToRead)
		reqLen, err = conn.Read(buf)
		if err != nil {
			return "", err
		}
		r += string(buf)
	}

	log.WithField("size", reqLen).Info("Got new request")
	return r, nil
}