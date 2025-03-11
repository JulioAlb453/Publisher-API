package repository

import (
	"encoding/json"
	"log"
	"publisher/core"
	"publisher/publisher/domain"
)

type EventRepository struct {
    broker *core.BrokerConnection
}

func NewEventRepository(broker *core.BrokerConnection) *EventRepository {
    return &EventRepository{broker: broker}
}

func (r *EventRepository) Publish(event domain.Event) error {
    eventJSON, err := json.Marshal(event)
    if err != nil {
        return err
    }
    log.Printf("Publicando evento en RabbitMQ: %s", eventJSON)
    return r.broker.Publish(eventJSON)
}