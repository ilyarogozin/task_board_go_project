package board

type BoardID struct {
	value string
}

func NewBoardID(value string) BoardID {
	return BoardID{value: value}
}

func (id BoardID) Value() string {
	return id.value
}