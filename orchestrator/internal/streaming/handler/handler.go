package handler

import (
	"context"

	"streaming/orchestrator/internal/streaming/model"
	"streaming/orchestrator/internal/streaming/pb"
)

type streamingController interface {
	Process(ctx context.Context, params model.QueryDto) (model.ResponseDto, error)
	CancelById(ctx context.Context, id int64) (bool, error)
	GetResult(ctx context.Context, id int64) (model.ResponseDto, error)
}

type Handler struct {
	pb.UnimplementedOrchestratorServer
	ctrl streamingController
}

func New(ctrl streamingController) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

func convertToPb(resp model.ResponseDto) *pb.Response {
	sources := make([]*pb.Source, 0, len(resp.Sources))
	for _, source := range resp.Sources {
		sources = append(sources, &pb.Source{
			Id:   source.Id,
			Link: source.Link,
		})
	}
	return &pb.Response{
		Id:      resp.Id,
		Status:  pb.ResponseStatus(resp.Status),
		Sources: sources,
	}
}
