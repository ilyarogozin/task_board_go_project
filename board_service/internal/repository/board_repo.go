package repository

import (
	"context"
	"encoding/json"
	"time"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BoardCreatedEvent struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id"`
}

type BoardRepository struct {
	db *pgxpool.Pool
}

func NewBoardRepository(db *pgxpool.Pool) *BoardRepository {
	return &BoardRepository{db: db}
}

func (r *BoardRepository) CreateBoard(
	ctx context.Context,
	title string,
	description string,
	ownerID string,
) (string, error) {

	if r.db == nil {
        return "", errors.New("pgx pool is nil (repository not initialized)")
    }

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
		boardID, title, description, ownerID,
	)
	if err != nil {
		return "", err
	}

	event := BoardCreatedEvent{
		ID:          boardID.String(),
		Title:       title,
		Description: description,
		OwnerID:     ownerID,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO outbox_events
		 (id, aggregate_id, event_type, payload, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(),
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