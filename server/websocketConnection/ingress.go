package websocketConnection

import (
	"fmt"
	"github.com/JoeReid/slb_websocket_server/server/schema"
	log "github.com/Sirupsen/logrus"
	"sync"
	"time"
)

func (c *Connection) ingressWorker(wg *sync.WaitGroup) {
	defer c.WS.Close()
	defer wg.Done()

	for {
		c.WS.SetReadDeadline(time.Now().Add(pongWait))

		var data schema.GenericJson

		err := c.jsonRead(&data)
		if err != nil {
			return // defer will Close conn
		}
		fmt.Println(data)
	}
}

func (c *Connection) jsonRead(dataPtr *schema.GenericJson) error {
	err := c.WS.ReadJSON(dataPtr)
	if err != nil {
		ip := c.WS.UnderlyingConn().RemoteAddr().String()

		log.WithFields(log.Fields{
			"ip":    ip,
			"error": err,
		}).Warn("Error reading from socket")
	}
	return err
}
