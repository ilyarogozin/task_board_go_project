package infra

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
)

func StartOutboxWorker(db *pgx.Conn) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "board-events",
	})

	go func() {
		for {
			rows, _ := db.Query(context.Background(),
				`DELETE FROM outbox_events
				 WHERE id IN (
					SELECT id FROM outbox_events
					LIMIT 10
					FOR UPDATE SKIP LOCKED
				 )
				 RETURNING type, payload`)

			for rows.Next() {
				var t string
				var p []byte
				rows.Scan(&t, &p)

				w.WriteMessages(context.Background(),
					kafka.Message{Value: p},
				)
			}
		}
	}()
	log.Println("Outbox worker started")
}