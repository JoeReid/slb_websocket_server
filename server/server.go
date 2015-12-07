package server

import (
	"fmt"
	"github.com/JoeReid/slb_websocket_server/listenerStopper"
	"github.com/JoeReid/slb_websocket_server/server/router"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"net"
	"sync"
)

func Run() {
	wg := sync.WaitGroup{}

	r := router.NewRouter()
	go r.Route()

	wg.Add(1)
	go startTCP(&wg, &r)

	wg.Wait()
}

func startTCP(wg *sync.WaitGroup, r *router.Router) {
	//get port from config
	tcpPort := fmt.Sprintf(
		":%v", viper.GetInt("server.tcp.port"),
	)

	l, err := net.Listen("tcp", tcpPort)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Panic("TCP Listen error")
	}

	l = listenerStopper.NewListener(l.(*net.TCPListener))
	defer l.Close()

	handle := tcpHandler{r}

	for {
		conn, err := l.Accept()
		log.WithFields(log.Fields{
			"ip": conn.RemoteAddr().String(),
		}).Debug("New tcp connection")

		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Panic("TCP Accept error")
		}
		go handle.Handle(conn)
	}

	wg.Done()
}
