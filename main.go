package main

import (
	"log"
	"time"

	"github.com/hdkshingala/labelwatcher/controller"
	localserver "github.com/hdkshingala/labelwatcher/server"
	"k8s.io/apiserver/pkg/server"
)

func main() {
	controller, err := controller.NewController()
	if err != nil {
		return
	}

	srv, err := localserver.NewServer(controller)
	if err != nil {
		return
	}

	stopChan := server.SetupSignalHandler()

	ch, err := srv.Config.SecureServingConfig.Serve(&srv.Mux, 30*time.Second, stopChan)
	if err != nil {
		log.Printf("Failed to create server. Error: %s\n", err.Error())
		return
	} else {
		<-ch
	}
}
