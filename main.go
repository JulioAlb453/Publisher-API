package main

import (
	"log"
	"publisher/core"
	infrastructure "publisher/publisher/infraestructure"
	"publisher/publisher/infraestructure/repository"
	"publisher/publisher/infraestructure/server"
)

func main() {
	broker, err := core.NewBrokerConnection()
	if err != nil {
		log.Fatal("Error al conectar con RabbitMQ: ", err)
	}

	defer func() {
		if broker != nil {
			broker.Close()
		}
	}()

	eventRepo := repository.NewEventRepository(broker)

	rabbitMQAdapter, err := infrastructure.NewRabbitMQAdapter("amqp://Julio:123456789@52.20.122.112:5672/", "Noticia Publicada")
	if err != nil {
		log.Fatal("Error al crear RabbitMQAdapter: ", err)
	}

	go func() {
		log.Println("Iniciando el servicio de escucha de eventos...")

		if err := rabbitMQAdapter.ListenToEvent(); err != nil {
			log.Fatalf("Error al escuchar eventos: %v", err)
		}
	}()

	httpServer := server.NewServer(eventRepo)

	if err := httpServer.Start(); err != nil {
		log.Fatal("Error al iniciar el servidor HTTP: ", err)
	}
}
