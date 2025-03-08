package repository

import (
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
	return r.broker.Publish(event.Payload)
}
