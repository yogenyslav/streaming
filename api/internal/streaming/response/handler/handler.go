package handler

import (
	"context"

	"streaming/api/internal/streaming/response/model"

	"github.com/go-playground/validator/v10"
)

type responseController interface {
	GetResult(ctx context.Context, queryId int64) (model.ResponseDto, error)
}

type Handler struct {
	ctrl      responseController
	validator *validator.Validate
}

func New(ctrl responseController) *Handler {
	return &Handler{
		ctrl:      ctrl,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}
