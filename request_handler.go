package go_redis_server

import (
	"net"

	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var repliesc chan Reply

type (
	Reply struct {
		id 		uuid.UUID
		message	string
	}

	Handler struct {
		Cfg			Redis
		Interrupt	*bool
	}
)

func (h *Handler) HandleRequest(conn net.Conn) {
	if repliesc == nil {
		repliesc = make(chan Reply, 0)
	}

	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	if err != nil {
		log.WithError(err).Error("Error reading request")
		conn.Write([]byte("-Error reading request\r\n"))
		conn.Close()
		return
	}

	log.WithField("size", reqLen).Info("Got new request")
	id := uuid.Must(uuid.NewV4())
	parsec<-Request{id, string(buf)}
	r := <-repliesc
	conn.Write([]byte(r.message))
	conn.Close()
}

func ReplyMessage(r string, id uuid.UUID) {
	repliesc <- Reply{id, r}
}

func ReplyError(err error, id uuid.UUID) {
	repliesc <- Reply{id, "-" + err.Error() + "\r\n"}
}