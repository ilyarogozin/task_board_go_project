package board

import (
	"errors"

	"github.com/google/uuid"
)

type OwnerID struct {
	value uuid.UUID
}

func NewOwnerID(value string) (OwnerID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return OwnerID{}, errors.New("invalid owner id")
	}
	return OwnerID{value: id}, nil
}

func (o OwnerID) Value() uuid.UUID {
	return o.value
}