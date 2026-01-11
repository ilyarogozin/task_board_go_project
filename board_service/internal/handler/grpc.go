package handler

import (
	"fmt"
	"context"

	board "github.com/ilyarogozin/task_board_go_project/gen/go/board"

	"board_service/internal/repository"
)

type BoardServer struct {
	board.UnimplementedBoardServiceServer
	repo repository.BoardRepository
}

func (s *BoardServer) CreateBoard(
	ctx context.Context,
	req *board.CreateBoardRequest,
) (*board.BoardResponse, error) {

	fmt.Println("=== CreateBoard received ===")
	fmt.Println("Title:", req.Title)
	fmt.Println("Description:", req.Description)
	fmt.Println("OwnerID:", req.OwnerId)
	fmt.Println("============================")

	id, err := s.repo.CreateBoardWithOutbox(
		ctx,
		req.Title,
		req.Description,
		req.OwnerId,
	)
	if err != nil {
		return nil, err
	}

	return &board.BoardResponse{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
		OwnerId:     req.OwnerId,
	}, nil
}