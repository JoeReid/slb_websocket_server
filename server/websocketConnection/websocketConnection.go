package websocketConnection

import (
	"fmt"
	"github.com/JoeReid/slb_websocket_server/server/router"
	"github.com/JoeReid/slb_websocket_server/server/schema"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 500 * time.Millisecond

	// Send pings to peer with this period. Must be less
	// than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 4096
)

type Connection struct {
	WS          *websocket.Conn
	EgressQueue chan schema.GenericJson
}

func (c *Connection) Regester(r *router.Router) {
	c.WS.SetReadDeadline(time.Now().Add(pongWait))

	var data schema.GenericJson

	err := c.jsonRead(&data)
	if err != nil {
		return // defer will Close conn
	}
	fmt.Println(data)

}

func (c *Connection) addToEgressQueue(data schema.GenericJson) {
	c.EgressQueue <- data
}

func (c *Connection) Work() {
	var wg sync.WaitGroup

	c.WS.SetReadLimit(maxMessageSize)
	c.WS.SetPongHandler(func(string) error {
		c.WS.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	wg.Add(2)
	go c.ingressWorker(&wg)
	go c.egressWorker(&wg)

	wg.Wait()
}
