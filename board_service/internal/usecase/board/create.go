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

	if err := ctx.Err(); err != nil {
		return "", err
	}

	title, err := board.NewBoardTitle(rawTitle)
	if err != nil {
		return "", err
	}

	ownerID, err := board.NewOwnerID(rawOwnerID)
	if err != nil {
		return "", err
	}

	description, err := board.NewBoardDescription(rawDescription)
	if err != nil {
		return "", err
	}

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