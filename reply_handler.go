package go_redis_server

import (
	"errors"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"net"
)

var repliesc chan Reply

type (
	Reply struct {
		id 		uuid.UUID
		message	string
	}

	ReplyHandle struct {
		Clients *map[uuid.UUID]net.Conn
		Worker
	}
)


func (r *ReplyHandle) Start() error {
	if r.Clients == nil {
		return	errors.New("reply handler can't start without given clients")
	}

	if repliesc == nil {
		repliesc = make(chan Reply, 0)
	}

	for i := 0; i < r.Cfg.ReplyWorkers; i++ {
		go r.run()
	}
	return nil
}

func (r *ReplyHandle) run() {
	for !*r.Interrupt {
		rep := <- repliesc
		c, ok := (*r.Clients)[rep.id]
		if ok {
			HandleReply(rep.message, c)
		} else {
			log.Error("got id without connection")
		}
	}
}

func HandleReply(msg string, conn net.Conn){
	conn.Write([]byte(msg))
}

func ReplyMessage(r string, id uuid.UUID) {
	repliesc <- Reply{id, r}
}

func ReplyError(err error, id uuid.UUID) {
	repliesc <- Reply{id, "-" + err.Error() + "\r\n"}
}

func ReplyNULL(id uuid.UUID) {
	repliesc <- Reply{id, buildRespNullBulkString()}
}