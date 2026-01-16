package board

type BoardDescription struct {
	value string
}

func NewBoardDescription(value string) BoardDescription {
	return BoardDescription{value: value}
}

func (d BoardDescription) Value() string {
	return d.value
}