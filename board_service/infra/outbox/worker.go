package outbox

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type Worker struct {
	db           *pgxpool.Pool
	writer       *kafka.Writer
	batchSize    int
	pollInterval time.Duration
}

type outboxEvent struct {
	id          string
	aggregateID string
	eventType   string
	payload     []byte
}

func NewOutboxWorker(
	db *pgxpool.Pool,
	writer *kafka.Writer,
) *Worker {
	return &Worker{
		db:           db,
		writer:       writer,
		batchSize:    10,
		pollInterval: 500 * time.Millisecond,
	}
}

func (w *Worker) Start(ctx context.Context) {
	go func() {
		log.Info().Msg("outbox worker started")

		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("outbox worker stopped")
				return
			default:
			}

			events, err := w.fetchAndMarkProcessing(ctx)
			if err != nil {
				log.Error().Err(err).Msg("failed to fetch outbox events")
				time.Sleep(time.Second)
				continue
			}

			for _, e := range events {
				if err := w.publishEvent(ctx, e); err != nil {
					log.Error().Err(err).
						Str("event_id", e.id).
						Msg("failed to publish outbox event")
					continue
				}

				if err := w.markPublished(ctx, e.id); err != nil {
					log.Error().Err(err).
						Str("event_id", e.id).
						Msg("failed to mark event as published")
					continue
				}

				log.Info().
				Str("event_id", e.id).
				Str("aggregate_id", e.aggregateID).
				Str("event_type", e.eventType).
				Msg("outbox event published")
			}

			time.Sleep(w.pollInterval)
		}
	}()
}

func (w *Worker) fetchAndMarkProcessing(
	ctx context.Context,
) ([]outboxEvent, error) {

	tx, err := w.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
		SELECT id, aggregate_id, event_type, payload
		FROM outbox_events
		WHERE status = 'pending'
		ORDER BY created_at
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`, w.batchSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]outboxEvent, 0, w.batchSize)

	for rows.Next() {
		var e outboxEvent
		if err := rows.Scan(
			&e.id,
			&e.aggregateID,
			&e.eventType,
			&e.payload,
		); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	if len(events) == 0 {
		return nil, tx.Commit(ctx)
	}

	ids := make([]string, 0, len(events))
	for _, e := range events {
		ids = append(ids, e.id)
	}

	_, err = tx.Exec(ctx, `
		UPDATE outbox_events
		SET status = 'processing'
		WHERE id = ANY($1)
	`, ids)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return events, nil
}

func (w *Worker) publishEvent(
	ctx context.Context,
	e outboxEvent,
) error {

	return w.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(e.aggregateID),
		Value: e.payload,
		Headers: []kafka.Header{
			{
				Key:   "event_type",
				Value: []byte(e.eventType),
			},
		},
	})
}

func (w *Worker) markPublished(
	ctx context.Context,
	eventID string,
) error {

	_, err := w.db.Exec(ctx, `
		UPDATE outbox_events
		SET status = 'published',
		    processed_at = now()
		WHERE id = $1
	`, eventID)

	return err
}