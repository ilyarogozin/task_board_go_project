package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Board struct {
	ID          uuid.UUID
	Title       string
	Description string
	OwnerID     uuid.UUID
	Columns     []*Column
}

func NewBoard(title, description string, ownerID uuid.UUID) (*Board, error) {
	if title == "" {
		return nil, errors.New("board title cannot be empty")
	}
	return &Board{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		OwnerID:     ownerID,
		Columns:     []*Column{},
	}, nil
}

func (b *Board) GetBoard() *Board {
	return b
}