package server

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type wsHandler struct {
}

func (wsh wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	// Check for error and bail early if found
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to upgrade HTTPHandler")
		return
	}

	c := &connection{
		ws:          ws,
		egressQueue: make(chan genericJson, 256),
	}
	c.regester()
	c.work()
}
