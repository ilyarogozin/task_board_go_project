package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"notification_service/internal/config"
	"notification_service/internal/infra/kafka"
	"notification_service/internal/infra/sender"
	"notification_service/internal/repository"
	"notification_service/internal/usecase"
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

	go func() {
		if err := consumer.Run(ctx); err != nil {
			log.Error().Err(err).Msg("kafka consumer exited with error")
			stop()
		}
	}()

	<-ctx.Done()
	log.Info().Msg("notification service shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := consumer.Close(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("failed to close kafka consumer gracefully")
	}

	log.Info().Msg("notification service stopped")
}