package websocketConnection

import (
	"fmt"
	"github.com/JoeReid/slb_websocket_server/server/router"
	"github.com/JoeReid/slb_websocket_server/server/schema"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

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
}

func (c *Connection) addToEgressQueue(data schema.GenericJson) {
	c.EgressQueue <- data
}

func (c *Connection) Work() {
	var wg sync.WaitGroup

	wg.Add(1)
	go c.ingressWorker(&wg)

	wg.Add(1)
	go c.egressWorker(&wg)

	wg.Wait()
}

func (c *Connection) ingressWorker(wg *sync.WaitGroup) {
	defer c.WS.Close()
	defer wg.Done()

	for {
		var data schema.GenericJson

		c.WS.SetReadLimit(maxMessageSize)
		c.WS.SetReadDeadline(time.Now().Add(pongWait))

		c.WS.SetPongHandler(func(string) error {
			c.WS.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		err := c.WS.ReadJSON(&data)
		if err != nil {
			ip := c.WS.UnderlyingConn().RemoteAddr().String()

			log.WithFields(log.Fields{
				"ip":    ip,
				"error": err,
			}).Warn("Error reading from socket: possible dead connection")

			return // defer will Close conn
		}
		fmt.Println(data)
	}
}

func (c *Connection) egressWorker(wg *sync.WaitGroup) {
	defer c.WS.Close()
	defer wg.Done()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		// If there is data to write
		case data, ok := <-c.EgressQueue:
			if !ok {
				c.closeWebsocket()
				return // defer will Close conn
			}

			// successful queue read
			ok = c.doWrite(data)
			if !ok {
				return // defer will Close conn
			}

		// If it is time to send a ping
		case <-ticker.C:
			ok := c.doPing()
			if !ok {
				return // defer will Close conn
			}
		}
	}
}

func (c *Connection) closeWebsocket() {
	ip := c.WS.UnderlyingConn().RemoteAddr().String()

	log.WithFields(log.Fields{
		"ip": ip,
	}).Warn("Error EgressQueue closed: closing connection")

	c.WS.SetWriteDeadline(time.Now().Add(writeWait))
	c.WS.WriteMessage(websocket.CloseMessage, []byte{})
}

func (c *Connection) doWrite(data schema.GenericJson) bool {
	c.WS.SetWriteDeadline(time.Now().Add(writeWait))

	err := c.WS.WriteJSON(data)
	if err != nil {
		ip := c.WS.UnderlyingConn().RemoteAddr().String()

		log.WithFields(log.Fields{
			"ip":    ip,
			"error": err,
		}).Warn("Error writing to socket: possible dead connection")

		return false
	}
	return true
}

func (c *Connection) doPing() bool {
	c.WS.SetWriteDeadline(time.Now().Add(writeWait))

	err := c.WS.WriteMessage(websocket.PingMessage, []byte{})
	if err != nil {
		ip := c.WS.UnderlyingConn().RemoteAddr().String()

		log.WithFields(log.Fields{
			"ip":    ip,
			"error": err,
		}).Warn("Error pinging socket: possible dead connection")

		return false
	}
	return true
}
