package server

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"sync"
)

type router struct {
	mutex         *sync.Mutex
	messageGroups map[string]*connectionPool
	queue         chan genericJson
}

func (r *router) route(data map[string]*json.RawMessage) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	g, ok := data["group"]
	if !ok {
		log.Error("No group specified in json")
		return //dont route
	}

	var groupKey string

	err := json.Unmarshal(*g, &groupKey)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Bad json formatting in router")
		return //dont route
	}

	pool, ok := r.messageGroups[groupKey]
	if !ok {
		log.WithFields(log.Fields{
			"group": groupKey,
		}).Error("Message group does not exist")
		return //dont route
	}

	pool.send(data)
}

func (r *router) work() {
	for {
		j := <-r.queue
		if j.Action != "message" {
			log.WithFields(log.Fields{
				"action": j.Action,
			}).Error("Incorrect action in router")
		} else {
			for _, elm := range j.Data {
				r.route(elm)
			}
		}
	}
}

func (r *router) subscribe(group string, conn *connection) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	pool, present := r.messageGroups[group]
	if present {
		pool.add(conn)
	} else {
		r.messageGroups[group] = &connectionPool{
			connections: *new([]*connection),
			mutex:       &sync.Mutex{},
		}
	}
}

type connectionPool struct {
	connections []*connection
	mutex       *sync.Mutex
}

func (c *connectionPool) send(data map[string]*json.RawMessage) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, conn := range c.connections {
		conn.egressQueue <- genericJson{
			Action: "message",
			Data: []map[string]*json.RawMessage{
				data,
			},
		}
	}
}

func (c *connectionPool) add(conn *connection) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.connections = append(c.connections, conn)
}
