package router

import (
	"github.com/JoeReid/slb_websocket_server/server/schema"
)

type Connection interface {
	Regester(*Router)
	Work()
	addToEgressQueue(schema.GenericJson)
}
