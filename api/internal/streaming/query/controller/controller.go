package controller

import (
	"context"

	"streaming/api/internal/shared"
	"streaming/api/internal/streaming/pb"
	"streaming/api/internal/streaming/query/model"
	respmodel "streaming/api/internal/streaming/response/model"

	"github.com/yogenyslav/pkg/storage/minios3"
	"go.opentelemetry.io/otel/trace"
)

type queryRepo interface {
	InsertOne(ctx context.Context, params model.QueryDao) (int64, error)
}

type responseRepo interface {
	InsertOne(ctx context.Context, params respmodel.ResponseDao) error
	UpdateOne(ctx context.Context, id int64, status shared.ResponseStatus) error
}

type Controller struct {
	qr           queryRepo
	rr           responseRepo
	orchestrator pb.OrchestratorClient
	s3           minios3.S3
	tracer       trace.Tracer
}

func New(qr queryRepo, rr responseRepo, orchestrator pb.OrchestratorClient, s3 minios3.S3, tracer trace.Tracer) *Controller {
	return &Controller{
		qr:           qr,
		rr:           rr,
		orchestrator: orchestrator,
		s3:           s3,
		tracer:       tracer,
	}
}
