package cmd

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "board_service/proto/board"
	"board_service/internal/handler"
	"board_service/infra"
	"board_service/internal/repository"
	"board_service/internal/usecase"
)

func main() {
	db := infra.MustPostgres()
	infra.StartOutboxWorker(db)

	uc := &usecase.BoardUsecase{
		DB:     db,
		Boards: repository.NewBoardRepo(db),
		Outbox: repository.NewOutboxRepo(),
	}

	l, _ := net.Listen("tcp", ":50051")
	s := grpc.NewServer()
	pb.RegisterBoardServiceServer(s, &handler.GRPC{UC: uc})

	log.Fatal(s.Serve(l))
}