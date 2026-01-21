package board

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrNilRepository = errors.New("board.Repository is nil")

type Repository interface {
	CreateBoard(
		ctx context.Context,
		title string,
		description string,
		ownerID uuid.UUID,
	) (string, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) (*Service, error) {
	if repo == nil {
		return nil, ErrNilRepository
	}

	return &Service{
		repo: repo,
	}, nil
}