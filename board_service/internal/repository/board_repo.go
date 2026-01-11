package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"board_service/internal/domain/board"
)

type BoardRepository struct {
	db *pgxpool.Pool
}

func NewBoardRepository(db *pgxpool.Pool) *BoardRepository {
	return &BoardRepository{db: db}
}

func (r *BoardRepository) CreateBoardWithOutbox(
	ctx context.Context,
	title string,
	description string,
	ownerId string,
) (string, error) {

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	boardID := uuid.New()
	now := time.Now()

	_, err = tx.Exec(
		ctx,
		`INSERT INTO boards (id, title, description, owner_id)
		 VALUES ($1, $2, $3, $4)`,
		boardID, title, description, ownerId,
	)
	if err != nil {
		return "", err
	}

	event := board.BoardCreatedEvent{
		Id:          boardID.String(),
		Title:       title,
		Description: description,
		OwnerId:     ownerId,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO outbox
		 (id, aggregate_type, aggregate_id, event_type, payload, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.New(),
		"board",
		boardID,
		"BoardCreated",
		payload,
		now,
	)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return boardID.String(), nil
}