package controller

import (
	"context"

	"streaming/api/internal/streaming/pb"
	"streaming/api/internal/streaming/response/model"

	"go.opentelemetry.io/otel/trace"
)

type responseRepo interface {
	FindOne(ctx context.Context, id int64) (model.ResponseDao, error)
}

type Controller struct {
	repo         responseRepo
	orchestrator pb.OrchestratorClient
	tracer       trace.Tracer
}

func New(repo responseRepo, orchestrator pb.OrchestratorClient, tracer trace.Tracer) *Controller {
	return &Controller{
		repo:         repo,
		orchestrator: orchestrator,
		tracer:       tracer,
	}
}
