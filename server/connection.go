package server

import (
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

type connection struct {
	ws          *websocket.Conn
	egressQueue chan genericJson
}

func (c *connection) regester() {
}

func (c *connection) work() {
	var wg sync.WaitGroup

	wg.Add(1)
	go c.ingressWorker(&wg)

	wg.Add(1)
	go c.egressWorker(&wg)

	wg.Wait()
}

func (c *connection) ingressWorker(wg *sync.WaitGroup) {
	defer c.ws.Close()
	defer wg.Done()

	for {
		var data genericJson

		c.ws.SetReadLimit(maxMessageSize)
		c.ws.SetReadDeadline(time.Now().Add(pongWait))

		c.ws.SetPongHandler(func(string) error {
			c.ws.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		err := c.ws.ReadJSON(&data)
		if err != nil {
			ip := c.ws.UnderlyingConn().RemoteAddr().String()

			log.WithFields(log.Fields{
				"ip":    ip,
				"error": err,
			}).Warn("Error reading from socket: possible dead connection")

			return // defer will Close conn
		}
	}
}

func (c *connection) egressWorker(wg *sync.WaitGroup) {
	defer c.ws.Close()
	defer wg.Done()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		// If there is data to write
		case data, ok := <-c.egressQueue:
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

func (c *connection) closeWebsocket() {
	ip := c.ws.UnderlyingConn().RemoteAddr().String()

	log.WithFields(log.Fields{
		"ip": ip,
	}).Warn("Error egressQueue closed: closing connection")

	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	c.ws.WriteMessage(websocket.CloseMessage, []byte{})
}

func (c *connection) doWrite(data genericJson) bool {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))

	err := c.ws.WriteJSON(data)
	if err != nil {
		ip := c.ws.UnderlyingConn().RemoteAddr().String()

		log.WithFields(log.Fields{
			"ip":    ip,
			"error": err,
		}).Warn("Error writing to socket: possible dead connection")

		return false
	}
	return true
}

func (c *connection) doPing() bool {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))

	err := c.ws.WriteMessage(websocket.PingMessage, []byte{})
	if err != nil {
		ip := c.ws.UnderlyingConn().RemoteAddr().String()

		log.WithFields(log.Fields{
			"ip":    ip,
			"error": err,
		}).Warn("Error pinging socket: possible dead connection")

		return false
	}
	return true
}
