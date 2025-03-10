package repository

import (
	"log"
	"publisher/core"
	"publisher/publisher/domain"
	infrastructure "publisher/publisher/infraestructure"
)

type EventRepository struct {
	broker *core.BrokerConnection
}

func NewEventRepository(broker *core.BrokerConnection) *EventRepository {
	return &EventRepository{broker: broker}
}

func (r *EventRepository) Publish(event domain.Event) error {
	return r.broker.Publish(event.Payload)
}

func (r *EventRepository) ListenToEvent() error {
	rabbitMQAdapter, err := infrastructure.NewRabbitMQAdapter("amqp://Julio:123456789@52.20.122.112:5672/", "newsQueue")
	if err != nil {
		return err
	}

	log.Println("Escuchando eventos...")
	return rabbitMQAdapter.ListenToEvent()
}