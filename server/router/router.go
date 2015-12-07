package router

import (
	"encoding/json"
	"github.com/JoeReid/slb_websocket_server/server/schema"
	log "github.com/Sirupsen/logrus"
	"sync"
)

var idCount int = 0
var idMu *sync.Mutex = &sync.Mutex{}

func NextID() int {
	idMu.Lock()
	defer idMu.Unlock()

	idCount++
	return idCount - 1
}

func NewRouter() Router {
	return Router{
		Queue:   make(chan schema.SingleMessage),
		Groups:  make(map[string]ConnectionPool),
		NoGroup: NewConnectionPool(),
	}
}

type Router struct {
	Queue   chan schema.SingleMessage
	Groups  map[string]ConnectionPool
	NoGroup ConnectionPool
}

func (r *Router) Subscribe(conn Connection, groups []string) {
	log.WithFields(log.Fields{
		"groups": groups,
	}).Debug("new subscriber")

	if len(groups) == 0 {
		r.NoGroup.AddConnection(conn)
		return
	}

	r.NoGroup.DeleteConnection(conn)

	for _, group := range groups {
		groupLogger := log.WithFields(log.Fields{
			"group": group,
		})

		pool, ok := r.Groups[group]
		if !ok {
			groupLogger.Info("creating new group")
			r.Groups[group] = NewConnectionPool()
			pool = r.Groups[group]
		}
		pool.AddConnection(conn)

	}
}

func (r *Router) Route() {
	for {
		msg := <-r.Queue
		r.handle(msg)
	}
}

func (r *Router) handle(msg schema.SingleMessage) {
	toMarshal := struct {
		Action string                        `json:"action"`
		Data   []map[string]*json.RawMessage `json:"data"`
	}{
		"message",
		[]map[string]*json.RawMessage{msg.WholeMessage},
	}

	marshaled, err := json.Marshal(toMarshal)
	if err != nil {
		log.WithFields(log.Fields{
			"message": msg,
		}).Error("failed to marshal message in router handler")
	}

	//send to the unsubscribed channels
	go r.NoGroup.Send(marshaled, msg.ConnectionID)

	if msg.Group == "" {
		// send to all
		for _, v := range r.Groups {
			go v.Send(marshaled, msg.ConnectionID)
		}
		return
	}

	cp, ok := r.Groups[msg.Group]
	if !ok {
		log.WithFields(log.Fields{
			"group": msg.Group,
		}).Error("mesage group does not exist")
		return
	}

	go cp.Send(marshaled, msg.ConnectionID)
}
