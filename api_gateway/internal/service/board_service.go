package service

import (
	"context"
	"time"

	boardpb "github.com/ilyarogozin/task_board_go_project/gen/go/board"
)

type BoardGRPC interface {
	CreateBoard(context.Context, *boardpb.CreateBoardRequest) (*boardpb.BoardResponse, error)
}

type BoardService struct {
	client BoardGRPC
}

func NewBoardService(client BoardGRPC) *BoardService {
	return &BoardService{client: client}
}

func (s *BoardService) CreateBoard(
	ctx context.Context,
	title, description, ownerID string,
) (*boardpb.BoardResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return s.client.CreateBoard(ctx, &boardpb.CreateBoardRequest{
		Title:       title,
		Description: description,
		OwnerId:     ownerID,
	})
}