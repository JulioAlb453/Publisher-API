package main

import (
	"log"
	"publisher/core"
	"publisher/publisher/infraestructure/repository"
	"publisher/publisher/infraestructure/server"
)

func main() {
	broker, err := core.NewBrokerConnection()

	if err != nil {
		log.Fatal("Error al conectar con RabbitMQ", err)
	}

	defer func() {
		if broker != nil {
			broker.Close()
		}
	}()

	eventRepo := repository.NewEventRepository(broker)
	httpServer := server.NewServer(eventRepo)

	if err := httpServer.Start(); err != nil {
		log.Fatal("Error al iniciar el servidor", err)
	}
}
