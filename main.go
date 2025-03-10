package main

import (
	"log"
	"os"
	"publisher/core"
	infrastructure "publisher/publisher/infraestructure"
	"publisher/publisher/infraestructure/repository"
	"publisher/publisher/infraestructure/server"

)

func main() {
        queueName := "Noticia Publicada" 
        RABBITMQ_URL:= os.Getenv("RABBITMQ_URL")

        broker, err := core.NewBrokerConnection(queueName)
        if err != nil {
                log.Fatalf("Error al conectar con RabbitMQ: %v", err)
        }
        defer broker.Close()

        rabbitMQAdapter, err := infrastructure.NewRabbitMQAdapter(RABBITMQ_URL, queueName)
        if err != nil {
                log.Fatalf("Error al crear RabbitMQAdapter: %v", err)
        }
        defer rabbitMQAdapter.Close()

        go func() {
                log.Println("Iniciando el servicio de escucha de eventos...")
                if err := rabbitMQAdapter.ListenToEvent(); err != nil {
                        log.Fatalf("Error al escuchar eventos: %v", err)
                }
        }()

        eventRepo := repository.NewEventRepository(broker)

        httpServer := server.NewServer(eventRepo)
        if err := httpServer.Start(); err != nil {
                log.Fatalf("Error al iniciar el servidor HTTP: %v", err)
        }
}