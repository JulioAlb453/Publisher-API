package main

import (
    "log"
    "os"
    "publisher/core"
    infrastructure "publisher/publisher/infraestructure"
    "publisher/publisher/infraestructure/repository"
    "publisher/publisher/infraestructure/server"
    "github.com/joho/godotenv"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error al cargar el archivo .env: %v", err)
    }

    queueName := "Noticia Publicada"

    rabbitMQURL := os.Getenv("RABBITMQ_URL")
    if rabbitMQURL == "" {
        log.Fatal("La variable de entorno RABBITMQ_URL no está definida")
    }

    broker, err := core.NewBrokerConnection(queueName)
    if err != nil {
        log.Fatalf("Error al conectar con RabbitMQ: %v", err)
    }
    defer broker.Close()

    rabbitMQAdapter, err := infrastructure.NewRabbitMQAdapter(rabbitMQURL, queueName)
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