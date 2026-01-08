package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type OutboxRepository interface {
	Add(ctx context.Context, tx pgx.Tx, eventType string, payload []byte) error
}

type outboxRepo struct{}

func NewOutboxRepo() OutboxRepository {
	return &outboxRepo{}
}

func (o *outboxRepo) Add(ctx context.Context, tx pgx.Tx, t string, p []byte) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO outbox_events (type, payload)
		 VALUES ($1,$2)`,
		t, p,
	)
	return err
}