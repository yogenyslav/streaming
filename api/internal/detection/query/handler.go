package query

import (
	"context"
	"net/http"
	"strconv"

	"streaming/api/internal/detection/query/model"
	"streaming/api/internal/shared"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type queryController interface {
	InsertOne(ctx context.Context, params model.QueryCreateReq) (int64, error)
	CancelById(ctx context.Context, id int64) error
}

type Handler struct {
	controller queryController
	validator  *validator.Validate
}

func NewHandler(controller queryController) *Handler {
	return &Handler{
		controller: controller,
		validator:  validator.New(validator.WithPrivateFieldValidation()),
	}
}

func (h *Handler) Video(ctx *fiber.Ctx) error {
	var (
		err     error
		queryId int64
	)

	rawFile, err := ctx.FormFile("source")
	if err != nil {
		return err
	}

	file, err := rawFile.Open()
	if err != nil {
		return err
	}

	req := model.QueryCreateReq{
		Video: file,
		Type:  shared.QueryTypeVideo,
		Name:  rawFile.Filename,
		Size:  rawFile.Size,
	}
	queryId, err = h.controller.InsertOne(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(model.QueryCreateResp{
		Id: queryId,
	})
}

func (h *Handler) Stream(ctx *fiber.Ctx) error {
	var (
		req     model.StreamCreateReq
		err     error
		queryId int64
	)

	if err = ctx.BodyParser(&req); err != nil {
		return err
	}
	if err = h.validator.Struct(&req); err != nil {
		return err
	}

	queryCreate := model.QueryCreateReq{
		Source: req.Source,
		Type:   shared.QueryTypeStream,
	}
	queryId, err = h.controller.InsertOne(ctx.Context(), queryCreate)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(model.QueryCreateResp{
		Id: queryId,
	})
}

func (h *Handler) CancelById(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	if err = h.controller.CancelById(ctx.Context(), id); err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}
