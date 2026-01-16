package board

import "errors"

type OwnerID struct {
	value string
}

func NewOwnerID(value string) (OwnerID, error) {
	if value == "" {
		return OwnerID{}, errors.New("owner id is required")
	}
	return OwnerID{value: value}, nil
}

func (id OwnerID) Value() string {
	return id.value
}