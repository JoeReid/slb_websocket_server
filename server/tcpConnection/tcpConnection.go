package tcpConnection

import (
	"encoding/json"
	"fmt"
	"github.com/JoeReid/slb_websocket_server/server/router"
	"github.com/JoeReid/slb_websocket_server/server/schema"
	"io/ioutil"
	"net"
)

type Connection struct {
	Conn        *net.Conn
	EgressQueue chan schema.GenericJson
}

func (c *Connection) Regester(r *router.Router) {
}

func (c *Connection) Work() {
	b, _ := ioutil.ReadAll(*c.Conn)

	var data schema.GenericJson
	json.Unmarshal(b, &data)
	fmt.Println(data)
	fmt.Println(string(b))
}

func (c *Connection) addToEgressQueue(data schema.GenericJson) {
	c.EgressQueue <- data
}
