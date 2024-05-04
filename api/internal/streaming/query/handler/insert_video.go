package handler

import (
	"net/http"

	"streaming/api/internal/shared"
	"streaming/api/internal/streaming/query/model"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/pkg/loctime"
)

func (h *Handler) InsertVideo(c *fiber.Ctx) error {
	rawFile, err := c.FormFile(shared.FormDataVideoKey)
	if err != nil {
		return err
	}

	file, err := rawFile.Open()
	if err != nil {
		return err
	}

	query := model.QueryCreate{
		Type:  shared.QueryTypeVideo,
		Name:  rawFile.Filename,
		Video: file,
		Size:  rawFile.Size,
		Ts:    loctime.GetLocalTime().Unix(),
	}
	id, err := h.ctrl.InsertOne(c.UserContext(), query)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(model.QueryCreateResp{Id: id})
}
