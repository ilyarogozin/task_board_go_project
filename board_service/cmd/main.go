package main

import (
	"context"
	"log"
	"net"
	"os"

	board "github.com/ilyarogozin/task_board_go_project/gen/go/board"
	"google.golang.org/grpc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"board_service/infra"
	"board_service/internal/handler"
)

func main() {
	err := godotenv.Load()
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is not set")
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to create pgx pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	defer pool.Close()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	board.RegisterBoardServiceServer(grpcServer, &handler.BoardServer{})

	go func() {
		log.Println("board_service gRPC listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("grpc serve error: %v", err)
		}
	}()

	writer := infra.NewKafkaWriter(
	[]string{"localhost:9092"},
	"board-events",
	)
	defer writer.Close()

	worker := infra.NewOutboxWorker(pool, writer)
	worker.Start(ctx)
}