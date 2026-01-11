package app

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	board "github.com/ilyarogozin/task_board_go_project/gen/go/board"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"board_service/infra"
	"board_service/internal/config"
	"board_service/internal/handler"
	"board_service/internal/repository"
)

func Run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error().Err(err).Msg("failed to load config")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Postgres
	pool, err := pgxpool.New(ctx, cfg.Database.DSN)
	if err != nil {
		log.Error().Err(err).Msg("failed to create pgx pool")
	}
	if err := pool.Ping(ctx); err != nil {
		log.Error().Err(err).Msg("failed to ping database")
	}
	defer pool.Close()

	// gRPC Server
	lis, err := net.Listen("tcp", cfg.Server.GRPCPort)
	if err != nil {
		log.Error().Err(err).Msg("failed to listen grpc")
	}
	grpcServer := grpc.NewServer()
	boardRepo := repository.NewBoardRepository(pool)
	boardServer := handler.NewBoardServer(boardRepo)
	board.RegisterBoardServiceServer(grpcServer, boardServer)

	go func() {
		log.Info().
			Str("service", "board_service").
			Str("grpc_port", cfg.Server.GRPCPort).
			Msg("gRPC server listening")
		if err := grpcServer.Serve(lis); err != nil {
			log.Error().Err(err).Msg("grpc serve error")
		}
	}()

	// Kafka
	writer := infra.NewKafkaWriter(
	cfg.Kafka.Brokers,
	cfg.Kafka.Topic,
	)
	defer writer.Close()

	// Outbox worker
	worker := infra.NewOutboxWorker(pool, writer)
	worker.Start(ctx)

	// block until exit signal
	<-sigCh
	log.Info().Msg("shutdown signal received")

	cancel()
	grpcServer.GracefulStop()

	return nil
}