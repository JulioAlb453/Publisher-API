package core

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

type BrokerConnection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewBrokerConnection() (*BrokerConnection, error) {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		conn.Close()
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		"events",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	log.Print("Conectado a RabbotMQ")
	return &BrokerConnection{conn: conn, channel: ch}, nil
}

func (b *BrokerConnection) Publish(event string) error {
	return b.channel.Publish(
		"",
		"events",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
}

func (b *BrokerConnection) Close() {
	b.channel.Close()
	b.conn.Close()
	log.Print("Desconectado de RabbotMQ")
}
