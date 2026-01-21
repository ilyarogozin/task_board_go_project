package board

import (
	"errors"
	"strings"
)

type BoardTitle struct {
	value string
}

const (
	boardTitleMinLen = 3
	boardTitleMaxLen = 100
)

func NewBoardTitle(value string) (BoardTitle, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return BoardTitle{}, errors.New("board title is required")
	}

	if len(value) < boardTitleMinLen {
		return BoardTitle{}, errors.New("board title is too short")
	}

	if len(value) > boardTitleMaxLen {
		return BoardTitle{}, errors.New("board title is too long")
	}

	return BoardTitle{value: value}, nil
}

func (t BoardTitle) Value() string {
	return t.value
}