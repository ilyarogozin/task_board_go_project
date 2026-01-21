package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	uc "notification_service/internal/usecase"
)

type BoardCreatedPayload struct {
	BoardID     string `json:"board_id"`
	BoardTitle  string `json:"board_title"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id"`
}

type Consumer struct {
	reader  *kafka.Reader
	useCase *uc.NotificationUseCase
}

func NewConsumer(
	brokers []string,
	topic string,
	useCase *uc.NotificationUseCase,
) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		MinBytes: 1e3,
		MaxBytes: 10e6,
	})

	return &Consumer{
		reader:  reader,
		useCase: useCase,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	log.Info().Msg("kafka consumer started")

	backoff := time.Second

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("kafka consumer context cancelled")
			return nil
		default:
		}

		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			log.Error().
				Err(err).
				Msg("failed to read kafka message, retrying")

			time.Sleep(backoff)
			if backoff < 30*time.Second {
				backoff *= 2
			}
			continue
		}

		backoff = time.Second

		eventType, err := getEventType(msg.Headers)
		if err != nil {
			log.Error().Err(err).Msg("event_type header missing")
			continue
		}

		switch eventType {
		case "BoardCreated":
			if err := c.handleBoardCreated(ctx, msg.Value); err != nil {
				log.Error().
					Err(err).
					Msg("failed to handle BoardCreated event")
			}
		default:
			log.Warn().
				Str("event_type", eventType).
				Msg("unknown event type received")
		}
	}
}

func (c *Consumer) Close(ctx context.Context) error {
	log.Info().Msg("closing kafka consumer")
	return c.reader.Close()
}

func (c *Consumer) handleBoardCreated(
	ctx context.Context,
	payload []byte,
) error {
	var event BoardCreatedPayload
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	userID, err := uuid.Parse(event.OwnerID)
	if err != nil {
		return err
	}

	return c.useCase.HandleBoardCreated(
		ctx,
		userID,
		event.BoardTitle,
	)
}

func getEventType(headers []kafka.Header) (string, error) {
	for _, h := range headers {
		if h.Key == "event_type" {
			return string(h.Value), nil
		}
	}
	return "", errors.New("event_type header not found")
}