package websocketConnection

import (
	"github.com/JoeReid/slb_websocket_server/server/schema"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

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

func (c *Connection) closeWebsocket() {
	ip := c.WS.UnderlyingConn().RemoteAddr().String()

	log.WithFields(log.Fields{
		"ip": ip,
	}).Warn("Error EgressQueue closed: closing connection")

	c.WS.SetWriteDeadline(time.Now().Add(writeWait))
	c.WS.WriteMessage(websocket.CloseMessage, []byte{})
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
