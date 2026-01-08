package handler

import (
	"context"

	pb "board_service/proto/board"
	"board_service/internal/usecase"
)

type GRPC struct {
	pb.UnimplementedBoardServiceServer
	UC *usecase.BoardUsecase
}

func (g *GRPC) CreateBoard(ctx context.Context, r *pb.CreateBoardRequest) (*pb.BoardResponse, error) {
	b, err := g.UC.CreateBoard(ctx, r.Title, r.Description, r.OwnerId)
	if err != nil {
		return nil, err
	}

	return &pb.BoardResponse{
		Id: b.ID, Title: b.Title, Description: b.Description, OwnerId: b.OwnerID,
	}, nil
}