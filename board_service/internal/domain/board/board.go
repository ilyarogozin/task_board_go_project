package board

type Board struct {
	id          BoardID
	title       BoardTitle
	description BoardDescription
	ownerID     OwnerID
}

func NewBoard(
	title BoardTitle,
	description BoardDescription,
	ownerID OwnerID,
) (*Board, error) {
	return &Board{
		title:       title,
		description: description,
		ownerID:     ownerID,
	}, nil
}

func (b *Board) ID() BoardID {
	return b.id
}

func (b *Board) Title() BoardTitle {
	return b.title
}

func (b *Board) Description() BoardDescription {
	return b.description
}

func (b *Board) OwnerID() OwnerID {
	return b.ownerID
}