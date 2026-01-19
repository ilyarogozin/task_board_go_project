package validator

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrTitleRequired   = errors.New("title is required")
	ErrOwnerIDRequired = errors.New("owner_id is required")
	ErrOwnerIDInvalid  = errors.New("owner_id must be valid UUID")
)

func ValidateCreateBoard(title, ownerID string) error {
	if title == "" {
		return ErrTitleRequired
	}
	if ownerID == "" {
		return ErrOwnerIDRequired
	}
	if _, err := uuid.Parse(ownerID); err != nil {
		return ErrOwnerIDInvalid
	}
	return nil
}