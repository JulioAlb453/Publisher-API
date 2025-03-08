package domain

type EventPublisher interface {
	Publish(event Event)  error

	ListenToEvent() error
}