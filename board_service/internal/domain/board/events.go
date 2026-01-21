package board

import "github.com/google/uuid"

type BoardCreatedEvent struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	OwnerID     uuid.UUID `json:"owner_id"`
}

func (BoardCreatedEvent) EventType() string {
	return "BoardCreated"
}