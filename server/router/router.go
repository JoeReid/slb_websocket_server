package router

import (
	"encoding/json"
	"github.com/JoeReid/slb_websocket_server/server/schema"
	log "github.com/Sirupsen/logrus"
)

func NewRouter() Router {
	return Router{
		Queue:   make(chan schema.SingleMessage),
		Groups:  *new(map[string]ConnectionPool),
		NoGroup: NewConnectionPool(),
	}
}

type Router struct {
	Queue   chan schema.SingleMessage
	Groups  map[string]ConnectionPool
	NoGroup ConnectionPool
}

func (r *Router) Subscribe(conn Connection, group string) {
	if group == "" {
		r.NoGroup.AddConnection(conn)
		return
	}

	r.NoGroup.DeleteConnection(conn)
	_, ok := r.Groups[group]
	if !ok {
		r.Groups[group] = NewConnectionPool()
	}
	pool := r.Groups[group]
	pool.AddConnection(conn)
}

func (r *Router) Route() {
	for {
		msg := <-r.Queue
		r.handle(msg)
	}
}

func (r *Router) handle(msg schema.SingleMessage) {
	toMarshal := struct {
		Action string                      `json:"action"`
		Data   map[string]*json.RawMessage `json:"data"`
	}{
		"message",
		msg.WholeMessage,
	}

	marshaled, err := json.Marshal(toMarshal)
	if err != nil {
		log.WithFields(log.Fields{
			"message": msg,
		}).Error("failed to marshal message in router handler")
	}

	//send to the unsubscribed channels
	go r.NoGroup.Send(marshaled)

	if msg.Group == "" {
		// send to all
		for _, v := range r.Groups {
			go v.Send(marshaled)
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

	go cp.Send(marshaled)
}
