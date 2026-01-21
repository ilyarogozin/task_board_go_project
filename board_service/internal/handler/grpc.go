package handler

import (
	"context"
	"errors"

	board "github.com/ilyarogozin/task_board_go_project/gen/go/board"
	boarduc "board_service/internal/usecase/board"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	if err := ctx.Err(); err != nil {
		log.Warn().Err(err).Msg("request context already cancelled")
		return nil, mapContextError(err)
	}

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

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			log.Warn().Err(err).Msg("request cancelled during CreateBoard")
			return nil, mapContextError(err)
		}

		log.Error().Err(err).Msg("failed to create board")
		return nil, status.Error(codes.Internal, "failed to create board")
	}

	if err := ctx.Err(); err != nil {
		log.Warn().Err(err).Msg("request cancelled before response")
		return nil, mapContextError(err)
	}

	return &board.BoardResponse{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
		OwnerId:     req.OwnerId,
	}, nil
}

func mapContextError(err error) error {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "request deadline exceeded")
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "request cancelled by client")
	default:
		return status.Error(codes.Internal, "context error")
	}
}