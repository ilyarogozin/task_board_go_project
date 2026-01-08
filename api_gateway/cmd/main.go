package cmd

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"

	pb "board_service/proto/board"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewBoardServiceClient(conn)

	app := fiber.New()

	app.Post("/boards", func(c *fiber.Ctx) error {
		var req pb.CreateBoardRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		res, err := client.CreateBoard(c.Context(), &req)
		if err != nil {
			return err
		}
		return c.JSON(res)
	})

	app.Get("/boards/:id", func(c *fiber.Ctx) error {
		res, err := client.GetBoard(
			c.Context(),
			&pb.GetBoardRequest{Id: c.Params("id")},
		)
		if err != nil {
			return err
		}
		return c.JSON(res)
	})

	log.Fatal(app.Listen(":8080"))
}