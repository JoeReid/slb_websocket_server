package tcpConnection

import (
	"fmt"
	"github.com/JoeReid/slb_websocket_server/server/router"
	//"github.com/JoeReid/slb_websocket_server/server/schema"
	"bufio"
	"io"
	"net"
	"sync"
)

type Connection struct {
	Conn        *net.Conn
	EgressQueue chan []byte
	Router      *router.Router
}

func (c *Connection) Work() {
	c.Router.Subscribe(c, "")
	wg := sync.WaitGroup{}
	close := make(chan struct{})

	wg.Add(2)
	go c.egressWork(&wg, close)
	go c.ingressWork(&wg, close)

	wg.Wait()
}

func (c *Connection) egressWork(wg *sync.WaitGroup, close chan struct{}) {
	for {
		select {
		case <-close:
			break
		case out := <-c.EgressQueue:
			con := *c.Conn
			con.Write(out)
		}
	}
	wg.Done()
}

func (c *Connection) ingressWork(wg *sync.WaitGroup, close chan struct{}) {
	for {
		msg, err := bufio.NewReader(*c.Conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				close <- struct{}{}
				break
			}
		}
		fmt.Println(msg)
	}
	wg.Done()
}

func (c *Connection) Send(data []byte) {
	c.EgressQueue <- data
}
