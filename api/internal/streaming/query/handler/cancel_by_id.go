package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CancelById(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	if err = h.ctrl.CancelById(c.UserContext(), id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
