package handler

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	"notification_service/internal/repository"
)

func StartConsumer(repo repository.NotificationRepository) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "board-events",
		GroupID: "notifications",
	})

	for {
		m, _ := r.ReadMessage(context.Background())
		log.Printf("Email sent: %s\n", string(m.Value))
		repo.Save(context.Background(), string(m.Value))
	}
}