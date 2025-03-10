package infrastructure

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQAdapter struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewRabbitMQAdapter(url, queue string) (*RabbitMQAdapter, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
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
		ch.Close()
		conn.Close()
		return nil, err
	}

	log.Printf("Conectado a RabbitMQ y cola '%s' declarada", queue)
	return &RabbitMQAdapter{conn: conn, channel: ch, queue: queue}, nil
}

func (a *RabbitMQAdapter) Publish(message string) error {
	return a.channel.Publish(
		"",
		a.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

func (a *RabbitMQAdapter) ListenToEvent() error {
	msgs, err := a.channel.Consume(
		a.queue,
		"Noticia Publicada",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Error al consumir el mensaje:", err)
		return err
	}

	log.Println("Esperando mensajes...")

	for msg := range msgs {
		var data map[string]interface{}
		if err := json.Unmarshal(msg.Body, &data); err != nil {
			log.Printf("Error al deserializar el mensaje: %v", err)
			continue
		}
		log.Printf("Evento recibido: %v", data)
	}

	return nil
}

func (a *RabbitMQAdapter) Close() {
	if a.channel != nil {
		a.channel.Close()
	}
	if a.conn != nil {
		a.conn.Close()
	}
	log.Println("Conexión a RabbitMQ cerrada")
}
