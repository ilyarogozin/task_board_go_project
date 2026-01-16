package main

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	board "github.com/ilyarogozin/task_board_go_project/gen/go/board"
)

type CreateBoardHTTPRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id"`
}

func main() {
	initLogger(zerolog.InfoLevel)

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to connect to gRPC")
	}
	defer conn.Close()

	log.Info().
		Str("address", "localhost:50051").
		Msg("connected to gRPC server")

	boardClient := board.NewBoardServiceClient(conn)

	app := fiber.New()

	app.Post("/boards", func(c *fiber.Ctx) error {
		var req CreateBoardHTTPRequest

		if err := c.BodyParser(&req); err != nil {
			log.Warn().
				Err(err).
				Msg("invalid json in request body")
			return fiber.NewError(fiber.StatusBadRequest, "invalid json")
		}

		log.Info().
			Interface("request", req).
			Msg("received CreateBoard request")

		if req.Title == "" {
			return fiber.NewError(fiber.StatusBadRequest, "title is required")
		}

		if req.OwnerID == "" {
			log.Debug().
				Bytes("raw_body", c.Body()).
				Msg("owner_id is empty")
			return fiber.NewError(fiber.StatusBadRequest, "owner_id is required")
		}

		if _, err := uuid.Parse(req.OwnerID); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "owner_id must be UUID")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		resp, err := boardClient.CreateBoard(ctx, &board.CreateBoardRequest{
			Title:       req.Title,
			Description: req.Description,
			OwnerId:     req.OwnerID,
		})
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to create board with gRPC")
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		log.Info().
			Str("board_id", resp.Id).
			Msg("board successfully created")

		return c.JSON(fiber.Map{
			"id":          resp.Id,
			"title":       resp.Title,
			"description": resp.Description,
			"owner_id":    resp.OwnerId,
		})
	})

	log.Info().
		Str("addr", ":8080").
		Msg("api_gateway HTTP listening")

	if err := app.Listen(":8080"); err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to start HTTP server")
	}
}

func initLogger(level zerolog.Level) {
	zerolog.TimeFieldFormat = time.RFC3339

	writer := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}

	log.Logger = zerolog.New(writer).
		Level(level).
		With().
		Timestamp().
		Logger()

	zerolog.SetGlobalLevel(level)
}