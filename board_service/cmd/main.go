package main

import (
	"context"
	"fmt"
	"log"
	"net"

	boardv1 "github.com/ilyarogozin/task_board_go_project/gen/go/board"
	"google.golang.org/grpc"
)

type boardServer struct {
	boardv1.UnimplementedBoardServiceServer
}

func (s *boardServer) CreateBoard(
	ctx context.Context,
	req *boardv1.CreateBoardRequest,
) (*boardv1.CreateBoardResponse, error) {

	fmt.Println("=== CreateBoard received ===")
	fmt.Println("Title:", req.Title)
	fmt.Println("Description:", req.Description)
	fmt.Println("OwnerID:", req.OwnerId)
	fmt.Println("============================")

	return &boardv1.CreateBoardResponse{
		Id: "test-board-id",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	boardv1.RegisterBoardServiceServer(grpcServer, &boardServer{})

	log.Println("board_service gRPC listening on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}