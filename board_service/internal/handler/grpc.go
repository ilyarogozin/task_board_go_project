package handler

import (
	"context"

	board "github.com/ilyarogozin/task_board_go_project/gen/go/board"
	"github.com/rs/zerolog/log"

	boarduc "board_service/internal/usecase/board"
)

type BoardServer struct {
	board.UnimplementedBoardServiceServer
	usecase *boarduc.Service
}

func NewBoardServer(usecase *boarduc.Service) *BoardServer {
	return &BoardServer{usecase: usecase}
}

func (s *BoardServer) CreateBoard(
	ctx context.Context,
	req *board.CreateBoardRequest,
) (*board.BoardResponse, error) {

	log.Info().
		Str("title", req.Title).
		Msg("CreateBoard request received")

	id, err := s.usecase.CreateBoard(
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