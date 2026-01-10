package main

import (
	"context"
	"fmt"
	"log"
	"net"

	board "github.com/ilyarogozin/task_board_go_project/gen/go/board"
	"google.golang.org/grpc"
)

type boardServer struct {
	board.UnimplementedBoardServiceServer
}

func (s *boardServer) CreateBoard(
	ctx context.Context,
	req *board.CreateBoardRequest,
) (*board.BoardResponse, error) {

	fmt.Println("=== CreateBoard received ===")
	fmt.Println("Title:", req.Title)
	fmt.Println("Description:", req.Description)
	fmt.Println("OwnerID:", req.OwnerId)
	fmt.Println("============================")

	return &board.BoardResponse{
		Id: "test-board-id",
		Title: req.Title,
		Description: req.Description,
		OwnerId: req.OwnerId,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	board.RegisterBoardServiceServer(grpcServer, &boardServer{})

	log.Println("board_service gRPC listening on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}