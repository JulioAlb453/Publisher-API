package application

import "publisher/publisher/domain"

type EventService struct {
	repo domain.EventPublisher
}

func NewEventService(repo domain.EventPublisher) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) PublishEvent(event domain.Event) error {
	return s.repo.Publish(event)
}