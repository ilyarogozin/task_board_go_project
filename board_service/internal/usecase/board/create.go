package board

import (
	"context"

	"board_service/internal/domain/board"
)

func (s *Service) CreateBoard(
	ctx context.Context,
	rawTitle string,
	rawDescription string,
	rawOwnerID string,
) (string, error) {

	title, err := board.NewBoardTitle(rawTitle)
	if err != nil {
		return "", err
	}

	ownerID, err := board.NewOwnerID(rawOwnerID)
	if err != nil {
		return "", err
	}

	description := board.NewBoardDescription(rawDescription)

	id, err := s.repo.CreateBoard(
		ctx,
		title.Value(),
		description.Value(),
		ownerID.Value(),
	)
	if err != nil {
		return "", err
	}

	return id, nil
}