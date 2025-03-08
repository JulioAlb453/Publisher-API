package application

import (
	"log"
	"publisher/publisher/domain"
)


type EventPubliserUseCase struct {
	Publisher domain.EventPublisher
}

func NewEventPubliserUseCase(publisher domain.EventPublisher) *EventPubliserUseCase {
    return &EventPubliserUseCase{Publisher: publisher}
}

func (uc *EventPubliserUseCase) ListenEvent() error {
	log.Print("Esperando eventos: ....")
	return uc.Publisher.ListenToEvent()
}

