package go_redis_server

import (
	log "github.com/sirupsen/logrus"
	"net"
)

func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	if err != nil {
		log.WithError(err).Error("Error reading request")
		conn.Write([]byte("Error reading request"))
	} else {
		log.WithField("size", reqLen).Info("Got new request")
		conn.Write([]byte("+PONG\r\n"))
	}
	conn.Close()
}