package handler

import (
	"net/http"

	"streaming/api/internal/shared"
	"streaming/api/internal/streaming/query/model"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) InsertStream(c *fiber.Ctx) error {
	var (
		req model.QueryStreamReq
		err error
	)

	if err = c.BodyParser(&req); err != nil {
		return err
	}
	if err = h.validator.Struct(&req); err != nil {
		return err
	}

	query := model.QueryCreate{
		Type: shared.QueryTypeStream,
		Name: req.Source,
	}
	id, err := h.ctrl.InsertOne(c.UserContext(), query)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(model.QueryCreateResp{Id: id})
}
