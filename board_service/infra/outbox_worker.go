package infra

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

type Worker struct {
	db     *pgxpool.Pool
	writer *kafka.Writer
	topic  string
}

func NewOutboxWorker(
	db *pgxpool.Pool,
	writer *kafka.Writer,
) *Worker {
	return &Worker{
		db:     db,
		writer: writer,
	}
}

func (w *Worker) Start(ctx context.Context) {
	go func() {
		for {
			rows, err := w.db.Query(ctx, `
				SELECT id, aggregate_id, event_type, payload
				FROM outbox_events
				WHERE processed IS FALSE
				ORDER BY created_at
				LIMIT 10
				FOR UPDATE SKIP LOCKED
			`)
			if err != nil {
				log.Println("outbox query error:", err)
				time.Sleep(time.Second)
				continue
			}

			for rows.Next() {
				var id string
				var aggregateID string
				var eventType string
				var payload []byte

				if err := rows.Scan(&id, &aggregateID, &eventType, &payload); err != nil {
					continue
				}

				err = w.writer.WriteMessages(ctx, kafka.Message{
					Key:   []byte(aggregateID),
					Value: payload,
					Headers: []kafka.Header{
						{Key: "event_type", Value: []byte(eventType)},
					},
				})
				if err != nil {
					log.Println("kafka write failed:", err)
					continue
				}

				_, _ = w.db.Exec(ctx,
					`UPDATE outbox_events SET processed = true WHERE id = $1`,
					id,
				)
			}

			rows.Close()
			time.Sleep(500 * time.Millisecond)
		}
	}()
}