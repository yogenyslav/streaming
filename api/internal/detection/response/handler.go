package response

import (
	"context"
	"net/http"
	"strconv"

	"streaming/api/internal/detection/response/model"

	"github.com/gofiber/fiber/v2"
)

type responseController interface {
	FindProcessed(ctx context.Context, queryId int64) (model.ResponseDto, error)
	GetStatic(ctx context.Context, name string) ([]byte, error)
}

type Handler struct {
	controller responseController
}

func NewHandler(controller responseController) *Handler {
	return &Handler{
		controller: controller,
	}
}

func (h *Handler) FindOneByQueryId(ctx *fiber.Ctx) error {
	queryId, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	resp, err := h.controller.FindProcessed(ctx.Context(), queryId)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

func (h *Handler) GetStatic(ctx *fiber.Ctx) error {
	name := ctx.Params("name")

	body, err := h.controller.GetStatic(ctx.Context(), name)
	if err != nil {
		return err
	}

	ctx.Response().Header.Add("Content-Type", "image/jpeg")
	return ctx.Status(http.StatusOK).Send(body)
}
