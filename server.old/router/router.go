package router

import (
	"encoding/json"
	"github.com/JoeReid/slb_websocket_server/server/schema"
	log "github.com/Sirupsen/logrus"
	"sync"
)

type Router struct {
	mutex         *sync.Mutex
	messageGroups map[string]*Pool
	Queue         chan schema.GenericJson
}

func NewRouter() *Router {
	return &Router{
		&sync.Mutex{},
		*new(map[string]*Pool),
		make(chan schema.GenericJson, 255),
	}
}

func (r *Router) route(data map[string]*json.RawMessage) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	g, ok := data["deviceid"]
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

func (r *Router) Work() {
	for {
		j := <-r.Queue
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

func (r *Router) Subscribe(group string, conn *Connection) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	pool, present := r.messageGroups[group]
	if present {
		pool.add(conn)
	} else {
		r.messageGroups[group] = &Pool{
			connections: *new([]*Connection),
			mutex:       &sync.Mutex{},
		}
	}
}
