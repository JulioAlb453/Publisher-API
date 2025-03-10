package infrastructure

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQAdapter struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewRabbitMQAdapter(rabbitMQURL, queue string) (*RabbitMQAdapter, error) {
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

	log.Println("Conectado a RabbitMQ correctamente")
	return &RabbitMQAdapter{
		conn:    conn,
		channel: ch,
		queue:   queue,
	}, nil
}

func (r *RabbitMQAdapter) Publish(event interface{}) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error al serializar el evento: %v", err)
		return err
	}

	log.Printf("Publicando mensaje en la cola %s: %s", r.queue, eventJSON) 
	return r.channel.Publish(
		"",
		r.queue,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, 
			ContentType:  "application/json",
			Body:         eventJSON,
		},
	)
}

func (r *RabbitMQAdapter) ListenToEvent() error {
	msgs, err := r.channel.Consume(
		r.queue,
		"",
		false, 
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Error al crear consumidor: %v", err)
		return err
	}

	log.Printf("Consumidor iniciado. Esperando mensajes en la cola %s...", r.queue) // Log agregado

	for msg := range msgs {
		log.Printf("Mensaje recibido (DeliveryTag: %d): %s", msg.DeliveryTag, msg.Body) // Log agregado

		var event map[string]interface{}
		err := json.Unmarshal(msg.Body, &event)
		if err != nil {
			log.Printf("Error al deserializar el mensaje (DeliveryTag: %d): %v", msg.DeliveryTag, err)
			msg.Nack(false, false) // Rechazar el mensaje y no reenviarlo
			continue
		}

		log.Printf("Evento recibido (DeliveryTag: %d): %+v", msg.DeliveryTag, event)

		// Procesar el mensaje...

		msg.Ack(false) // Confirmar el mensaje después de procesarlo
	}

	return nil
}

// Close cierra la conexión y el canal de RabbitMQ de manera segura.
func (r *RabbitMQAdapter) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	log.Println("Desconectado de RabbitMQ")
}
