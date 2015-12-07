package router

import (
	"errors"
)

var (
	ConnectionNotInPoolError     = errors.New("Connection not in pool")
	ConnectionAlreadyInPoolError = errors.New("Connection already in pool")
)

func NewConnectionPool() ConnectionPool {
	return ConnectionPool{
		make(map[Connection]bool),
	}
}

type ConnectionPool struct {
	Connections map[Connection]bool
}

func (c *ConnectionPool) AddConnection(conn Connection) error {
	_, exists := c.Connections[conn]
	if exists {
		return ConnectionAlreadyInPoolError
	}

	c.Connections[conn] = true
	return nil
}

func (c *ConnectionPool) DeleteConnection(conn Connection) error {
	_, exists := c.Connections[conn]
	if !exists {
		return ConnectionNotInPoolError
	}

	delete(c.Connections, conn)
	return nil
}

func (c *ConnectionPool) Send(msg []byte, excludeID int) {
	send := func(conn Connection) {
		if excludeID != conn.GetID() {
			conn.Send(msg)
		}
		// maybe cleanup closed conns here later
	}

	for k, _ := range c.Connections {
		go send(k)
	}
}
