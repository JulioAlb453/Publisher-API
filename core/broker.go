package core

import (
        "log"
        "os"

        "github.com/joho/godotenv"
        "github.com/streadway/amqp"
)

type BrokerConnection struct {
        conn    *amqp.Connection
        channel *amqp.Channel
        queue   string
}

func NewBrokerConnection(queue string) (*BrokerConnection, error) {
        err := godotenv.Load()
        if err != nil {
                log.Fatal("Error al cargar el archivo .env")
        }

        rabbitMQURL := os.Getenv("RABBITMQ_URL")
        if rabbitMQURL == "" {
                rabbitMQURL = "amqp://guest:guest@52.20.122.112:5672/"
        }
        conn, err := amqp.Dial(rabbitMQURL)
        if err != nil {
                log.Printf("Error al conectar con RabbitMQ: %v", err)
                return nil, err
        }

        ch, err := conn.Channel()
        if err != nil {
                log.Printf("Error al abrir canal en RabbitMQ: %v", err)
                conn.Close()
                return nil, err
        }

        _, err = ch.QueueDeclare(
                queue,
                true,
                false,
                false,
                false,
                nil,
        )
        if err != nil {
                log.Printf("Error al declarar la cola: %v", err)
                ch.Close()
                conn.Close()
                return nil, err
        }

        log.Printf("Conectado a RabbitMQ correctamente y cola '%s' declarada", queue)
        return &BrokerConnection{conn: conn, channel: ch, queue: queue}, nil
}

func (b *BrokerConnection) Publish(event []byte) error {
        if b.channel == nil {
                return amqp.ErrClosed
        }
        return b.channel.Publish(
                "",           
                b.queue,     
                false,         
                false,          
                amqp.Publishing{
                        ContentType: "application/json",
                        Body:        event,
                },
        )
}

func (b *BrokerConnection) Close() {
        if b.channel != nil {
                b.channel.Close()
        }
        if b.conn != nil {
                b.conn.Close()
        }
        log.Println("Desconectado de RabbitMQ")
}