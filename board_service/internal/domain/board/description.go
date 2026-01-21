package board

import (
	"strings"
	"errors"
)

type BoardDescription struct {
	value string
}

const boardDescriptionMaxLen = 1000

func NewBoardDescription(value string) (BoardDescription, error) {
	value = strings.TrimSpace(value)

	if len(value) > boardDescriptionMaxLen {
		return BoardDescription{}, errors.New("board description too long")
	}

	return BoardDescription{value: value}, nil
}

func (d BoardDescription) Value() string {
	return d.value
}