package server

import (
	"github.com/JoeReid/slb_websocket_server/server/router"
	"github.com/JoeReid/slb_websocket_server/server/schema"
	"github.com/JoeReid/slb_websocket_server/server/tcpConnection"
	"net"
)

type tcpHandler struct {
	Router *router.Router
}

func (t *tcpHandler) Handle(conn *net.Conn) {
	c := &tcpConnection.Connection{
		Conn:        conn,
		EgressQueue: make(chan schema.GenericJson, 256),
	}
	c.Regester(t.Router)
	c.Work()
}
