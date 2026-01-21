package app

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"api_gateway/internal/handler"
	"api_gateway/internal/service"
	grpcclient "api_gateway/internal/transport/grpc"
)

type App struct {
	fiber *fiber.App
	grpc  *grpcclient.BoardClient
}

func New() (*App, error) {
	grpcClient, err := grpcclient.NewBoardClient("localhost:50051")
	if err != nil {
		return nil, err
	}

	boardService := service.NewBoardService(grpcClient)
	boardHandler := handler.NewBoardHandler(boardService)

	app := fiber.New()
	boardHandler.Register(app)

	return &App{
		fiber: app,
		grpc:  grpcClient,
	}, nil
}

func (a *App) Run() error {
	return a.fiber.Listen(":8080")
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.fiber.Shutdown(); err != nil {
		return err
	}
	return a.grpc.Close()
}