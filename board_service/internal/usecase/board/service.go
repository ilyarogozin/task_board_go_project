package board

import "context"

type Repository interface {
	CreateBoard(
		ctx context.Context,
		title string,
		description string,
		ownerID string,
	) (string, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("board.Repository is nil")
	}
	return &Service{repo: repo}
}