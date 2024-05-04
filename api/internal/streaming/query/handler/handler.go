package handler

import (
	"context"

	"streaming/api/internal/streaming/query/model"

	"github.com/go-playground/validator/v10"
)

type queryController interface {
	InsertOne(ctx context.Context, params model.QueryCreate) (int64, error)
	CancelById(ctx context.Context, id int64) error
}

type Handler struct {
	ctrl      queryController
	validator *validator.Validate
}

func New(ctrl queryController) *Handler {
	return &Handler{
		ctrl:      ctrl,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}
