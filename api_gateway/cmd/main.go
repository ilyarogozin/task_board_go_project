package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"api_gateway/internal/app"
)

func main() {
	initLogger(zerolog.InfoLevel)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	application, err := app.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init app")
	}

	go func() {
		if err := application.Run(); err != nil {
			log.Fatal().Err(err).Msg("http server error")
		}
	}()

	log.Info().Msg("api gateway started")

	<-ctx.Done()
	log.Info().Msg("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("graceful shutdown failed")
	}

	log.Info().Msg("api gateway stopped")
}

func initLogger(level zerolog.Level) {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = zerolog.New(os.Stderr).
		Level(level).
		With().
		Timestamp().
		Logger()
}