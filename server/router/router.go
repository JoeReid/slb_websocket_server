package router

import (
	"github.com/JoeReid/slb_websocket_server/server/schema"
	log "github.com/Sirupsen/logrus"
)

type Router struct {
	Queue   chan schema.SingleMessage
	Groups  map[string]ConnectionPool
	NoGroup ConnectionPool
}

func (r *Router) Route() {
	for {
		msg := <-r.Queue
		r.handle(msg)
	}
}

func (r *Router) handle(msg schema.SingleMessage) {
	//send to the unsubscribed channels
	go r.NoGroup.Send(msg)

	if msg.Group == "" {
		// send to all
		for _, v := range r.Groups {
			go v.Send(msg)
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

	go cp.Send(msg)
}
