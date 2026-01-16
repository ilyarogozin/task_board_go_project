package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"notification_service/internal/usecase"
	"notification_service/internal/config"
	"notification_service/internal/infra/kafka"
	"notification_service/internal/repository"
	"notification_service/internal/infra/sender"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.Database.DSN)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect postgres")
	}
	defer pool.Close()

	repo := repository.NewNotificationRepository(pool)
	emailSender := sender.NewEmailSender()

	useCase := usecase.NewNotificationUseCase(repo, emailSender)

	consumer := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.Topic,
		useCase,
	)

	log.Info().Msg("notification service started")

	if err := consumer.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("kafka consumer stopped")
	}
}