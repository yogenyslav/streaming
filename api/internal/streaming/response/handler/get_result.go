package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetResult(c *fiber.Ctx) error {
	queryId, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	result, err := h.ctrl.GetResult(c.Context(), queryId)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)
}
