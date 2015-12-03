package server

import (
	"github.com/JoeReid/slb_websocket_server/server/router"
	"github.com/JoeReid/slb_websocket_server/server/schema"
	"github.com/JoeReid/slb_websocket_server/server/websocketConnection"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsHandler struct {
	Router *router.Router
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

	log.WithFields(log.Fields{
		"ip": ws.UnderlyingConn().RemoteAddr().String(),
	}).Debug("New websocket connection")

	c := &websocketConnection.Connection{
		WS:          ws,
		EgressQueue: make(chan schema.GenericJson, 256),
	}
	c.Regester(wsh.Router)
	c.Work()
}
