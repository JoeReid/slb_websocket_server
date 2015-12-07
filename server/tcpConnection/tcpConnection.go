package tcpConnection

import (
	"bufio"
	"encoding/json"
	"github.com/JoeReid/slb_websocket_server/server/router"
	"github.com/JoeReid/slb_websocket_server/server/schema"
	log "github.com/Sirupsen/logrus"
	"io"
	"net"
	"sync"
)

type Connection struct {
	Conn        net.Conn
	EgressQueue chan []byte
	Router      *router.Router
}

func (c *Connection) Work() {
	c.Router.Subscribe(c, []string{})
	wg := sync.WaitGroup{}
	close := make(chan struct{})

	wg.Add(2)
	go c.egressWork(&wg, close)
	go c.ingressWork(&wg, close)

	wg.Wait()
}

func (c *Connection) egressWork(wg *sync.WaitGroup, close chan struct{}) {
	for {
		select {
		case <-close:
			log.Debug("egress loop signaled to close")
			goto stop
		case out := <-c.EgressQueue:
			c.Conn.Write(out)
		}
	}

stop:
	log.Debug("Egress loop stopped")
	wg.Done()
}

func (c *Connection) ingressWork(wg *sync.WaitGroup, close chan struct{}) {
	for {
		msg, err := bufio.NewReader(c.Conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Debug("Closing conn due to EOF")
			} else {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Failed to read from conn... closing conn")
			}

			close <- struct{}{}
			break
		}

		log.WithFields(log.Fields{
			"message": msg,
		}).Debug("new message")

		c.handleMessage([]byte(msg), close)
	}
	wg.Done()
}

func (c *Connection) handleMessage(data []byte, close chan struct{}) {
	logger := log.WithFields(log.Fields{
		"json": string(data),
	})
	var generic schema.GenericJson

	err := json.Unmarshal(data, &generic)
	if err != nil {
		logger.Error("Malformed json in tcp handleMessage")
		return
	}

	switch {
	case generic.Action == "message":
		close <- struct{}{}
		msg, err := generic.ToMessage()
		if err != nil {
			logger.WithFields(log.Fields{
				"error": err,
			}).Error("Fail to cast to MessageJson")
			return
		}

		messageArray, err := msg.ToMessageArray()
		if err != nil {
			logger.WithFields(log.Fields{
				"error": err,
			}).Error("Fail to cast to MessageArray")
			return
		}

		for _, single := range messageArray {
			c.Router.Queue <- single
		}

	}
}

func (c *Connection) Send(data []byte) {
	c.EgressQueue <- data

	log.WithFields(log.Fields{
		"message": string(data),
	}).Debug("New message recieved from router")
}
