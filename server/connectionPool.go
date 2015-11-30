package server

import (
	"encoding/json"
	"sync"
)

type connectionPool struct {
	connections []*connection
	mutex       *sync.Mutex
}

func (c *connectionPool) send(data map[string]*json.RawMessage) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, conn := range c.connections {
		go func() {
			conn.egressQueue <- genericJson{
				Action: "message",
				Data: []map[string]*json.RawMessage{
					data,
				},
			}
		}()
	}
}

func (c *connectionPool) add(conn *connection) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.connections = append(c.connections, conn)
}
