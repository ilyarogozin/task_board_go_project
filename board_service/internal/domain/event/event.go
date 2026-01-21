package event

type DomainEvent interface {
	EventType() string
}