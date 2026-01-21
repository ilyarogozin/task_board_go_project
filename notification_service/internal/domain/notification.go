package domain

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Message   string
	CreatedAt time.Time
}

func NewNotification(userID uuid.UUID, message string) *Notification {
	return &Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Message:   message,
		CreatedAt: time.Now().UTC(),
	}
}