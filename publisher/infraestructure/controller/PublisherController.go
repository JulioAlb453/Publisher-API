package infrastructure

import (
	"log"
	"encoding/json"
	"github.com/streadway/amqp"
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

	return r.channel.Publish(
		"",              
		"Noticia Publicada",           
		false,         
		false,           
		amqp.Publishing{
			ContentType: "application/json", 
			Body:        eventJSON,          
		},
	)
}

func (r *RabbitMQAdapter) ListenToEvent() error {
	msgs, err := r.channel.Consume(
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
		return err
	}

	for msg := range msgs {
		log.Printf("Mensaje recibido: %s", msg.Body)

		var event map[string]interface{}
		err := json.Unmarshal(msg.Body, &event)
		if err != nil {
			log.Printf("Error al deserializar el mensaje: %v", err)
			continue
		}

		log.Printf("Evento recibido: %+v", event)

	}

	return nil
}
func (r *RabbitMQAdapter) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	log.Println("Desconectado de RabbitMQ")
}
