package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	),
	)
	if err != nil {
		log.Fatalf("failed to connect to gRPC: %v", err)
	}
	defer conn.Close()

	boardClient := board.NewBoardServiceClient(conn)

	app := fiber.New()

	app.Post("/boards", func(c *fiber.Ctx) error {
		var req CreateBoardHTTPRequest

		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid json")
		}
		fmt.Println("Received CreateBoard request:", req)

		if req.Title == "" {
			return fiber.NewError(fiber.StatusBadRequest, "title is required")
		}

		if req.OwnerID == "" {
			fmt.Printf("DEBUG RAW BODY: %s\n", c.Body())
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
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"id": resp.Id,
			"title": resp.Title,
			"description": resp.Description,
			"owner_id": resp.OwnerId,
		})
	})

	log.Println("api_gateway HTTP listening on :8080")
	log.Fatal(app.Listen(":8080"))
}