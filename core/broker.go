package core

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
	"github.com/joho/godotenv"
)

type BrokerConnection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewBrokerConnection() (*BrokerConnection, error) {
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
		"Noticia Publicada", 
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

	log.Println("Conectado a RabbitMQ correctamente")
	return &BrokerConnection{conn: conn, channel: ch}, nil
}

func (b *BrokerConnection) Publish(event string) error {
	if b.channel == nil {
		return amqp.ErrClosed
	}
	return b.channel.Publish(
		"",               
		"Noticia Publicada", 
		false,             
		false,             
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event), 
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

func (b *BrokerConnection) ListenToQueue() (<-chan amqp.Delivery, error) {
	if b.channel == nil {
		return nil, fmt.Errorf("canal no abierto")
	}

	msgs, err := b.channel.Consume(
		"Noticia Publicada", 
		"",                  
		true,               
		false,               
		false,              
		false,              
		nil,                 
	)
	if err != nil {
		log.Printf("Error al crear consumidor: %v", err)
		return nil, err
	}

	return msgs, nil
}
