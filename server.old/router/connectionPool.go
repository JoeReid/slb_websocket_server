package router

import (
	"encoding/json"
	"github.com/JoeReid/slb_websocket_server/server/schema"
	"sync"
)

type Pool struct {
	connections []*Connection
	mutex       *sync.Mutex
}

func (c *Pool) send(data map[string]*json.RawMessage) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, conn := range c.connections {
		go func() {
			deref := *conn
			deref.addToEgressQueue(schema.GenericJson{
				Action: "message",
				Data: []map[string]*json.RawMessage{
					data,
				},
			})
		}()
	}
}

func (c *Pool) add(conn *Connection) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.connections = append(c.connections, conn)
}
