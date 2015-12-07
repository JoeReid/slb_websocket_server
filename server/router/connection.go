package router

import (
//"github.com/JoeReid/slb_websocket_server/server/schema"
)

type Connection interface {
	Send(msg []byte)
	Work()
	GetID() int
}
