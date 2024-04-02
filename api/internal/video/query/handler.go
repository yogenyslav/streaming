package query

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"streaming/internal/video/query/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type queryController interface {
	InsertOne(ctx context.Context, processing map[int64]context.CancelFunc, params model.QueryCreateReq) (int64, error)
	CancelById(ctx context.Context, processing map[int64]context.CancelFunc, id int64) error
}

type Handler struct {
	controller queryController
	validator  *validator.Validate
	processing map[int64]context.CancelFunc
}

func NewHandler(controller queryController) *Handler {
	return &Handler{
		controller: controller,
		validator:  validator.New(validator.WithPrivateFieldValidation()),
		processing: make(map[int64]context.CancelFunc),
	}
}

func (h *Handler) File(ctx *fiber.Ctx) error {
	var (
		err     error
		queryId int64
	)

	file, err := ctx.FormFile("src")
	if err != nil {
		return err
	}

	fileNameSplit := strings.Split(file.Filename, ".")
	fileExt := "." + fileNameSplit[len(fileNameSplit)-1]
	source := "./static/" + uuid.New().String()

	if err = ctx.SaveFile(file, source+"_"+fileExt); err != nil {
		return err
	}

	req := model.QueryCreateReq{
		Source:  source,
		FileExt: fileExt,
	}
	queryId, err = h.controller.InsertOne(ctx.Context(), h.processing, req)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(model.QueryCreateResp{
		Id: queryId,
	})
}

func (h *Handler) Link(ctx *fiber.Ctx) error {
	var (
		req     model.QueryCreateReq
		err     error
		queryId int64
	)

	if err = ctx.BodyParser(&req); err != nil {
		return err
	}
	if err = h.validator.Struct(&req); err != nil {
		return err
	}

	queryId, err = h.controller.InsertOne(ctx.Context(), h.processing, req)
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

	if err = h.controller.CancelById(ctx.Context(), h.processing, id); err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}
