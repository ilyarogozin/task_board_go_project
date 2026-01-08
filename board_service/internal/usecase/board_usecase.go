package usecase

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"board_service/internal/domain"
	"board_service/internal/repository"
)

type BoardUsecase struct {
	DB     *pgx.Conn
	Boards repository.BoardRepository
	Outbox repository.OutboxRepository
}

func (uc *BoardUsecase) CreateBoard(ctx context.Context, title, desc string, ownerID uuid.UUID) (*domain.Board, error) {
	tx, err := uc.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	board := domain.Board{
		ID:          uuid.New(),
		Title:       title,
		Description: desc,
		OwnerID:     ownerID,
	}

	if err := uc.Boards.Create(ctx, tx, board); err != nil {
		return nil, err
	}

	payload, _ := json.Marshal(board)
	if err := uc.Outbox.Add(ctx, tx, "BoardCreated", payload); err != nil {
		return nil, err
	}

	return &board, tx.Commit(ctx)
}