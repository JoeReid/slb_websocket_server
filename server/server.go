package server

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func Run() {
	//get port from config
	port := fmt.Sprintf(":%v", viper.GetInt("server.port"))

	handle := wsHandler{}

	http.Handle("/", handle)

	log.Info("Starting websocket server")
	if err := http.ListenAndServe(port, nil); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Panic("ListenAndServe error")
	}
}
