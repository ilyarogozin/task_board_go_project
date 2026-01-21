package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"api_gateway/internal/service"
	"api_gateway/internal/validator"
)

type BoardHandler struct {
	service *service.BoardService
}

type CreateBoardHTTPRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id"`
}

func NewBoardHandler(service *service.BoardService) *BoardHandler {
	return &BoardHandler{service: service}
}

func (h *BoardHandler) Register(app *fiber.App) {
	app.Post("/boards", h.CreateBoard)
}

func (h *BoardHandler) CreateBoard(c *fiber.Ctx) error {
	var req CreateBoardHTTPRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}

	if err := validator.ValidateCreateBoard(req.Title, req.OwnerID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	resp, err := h.service.CreateBoard(
		c.Context(),
		req.Title,
		req.Description,
		req.OwnerID,
	)
	if err != nil {
		log.Error().Err(err).Msg("create board failed")
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.JSON(fiber.Map{
		"id":          resp.Id,
		"title":       resp.Title,
		"description": resp.Description,
		"owner_id":    resp.OwnerId,
	})
}