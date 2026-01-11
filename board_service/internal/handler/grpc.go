package handler

import (
	"context"

	board "github.com/ilyarogozin/task_board_go_project/gen/go/board"
	"github.com/rs/zerolog/log"

	"board_service/internal/repository"
)

type BoardServer struct {
	board.UnimplementedBoardServiceServer
	repo *repository.BoardRepository
}

func NewBoardServer(repo *repository.BoardRepository) *BoardServer {
    if repo == nil {
        panic("BoardRepository is nil")
    }
    return &BoardServer{repo: repo}
}

func (s *BoardServer) CreateBoard(
	ctx context.Context,
	req *board.CreateBoardRequest,
) (*board.BoardResponse, error) {

	log.Info().
		Str("title", req.Title).
		Str("description", req.Description).
		Str("owner_id", req.OwnerId).
		Msg("CreateBoard request received")

	id, err := s.repo.CreateBoardWithOutbox(
		ctx,
		req.Title,
		req.Description,
		req.OwnerId,
	)
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("board_id", id).
		Msg("Board created successfully")

	return &board.BoardResponse{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
		OwnerId:     req.OwnerId,
	}, nil
}