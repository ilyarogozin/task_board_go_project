package board

import "errors"

type BoardTitle struct {
	value string
}

func NewBoardTitle(value string) (BoardTitle, error) {
	if value == "" {
		return BoardTitle{}, errors.New("board title can't be empty")
	}
	if len(value) > 255 {
		return BoardTitle{}, errors.New("board title too long")
	}
	return BoardTitle{value: value}, nil
}

func (t BoardTitle) Value() string {
	return t.value
}