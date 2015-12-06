package tcpConnection

import (
	"fmt"
	"github.com/JoeReid/slb_websocket_server/server/router"
	"github.com/JoeReid/slb_websocket_server/server/schema"
	"io/ioutil"
	"net"
)

type Connection struct {
	Conn        *net.Conn
	EgressQueue chan schema.GenericJson
	router      *router.Router
}

func (c *Connection) Regester(r *router.Router) {
	c.router = r

	b, _ := ioutil.ReadAll(*c.Conn)
	fmt.Println(string(b))

}

func (c *Connection) Work() {
}

func (c *Connection) addToEgressQueue(data schema.GenericJson) {
	c.EgressQueue <- data
}
