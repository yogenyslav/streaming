package controller

import (
	"context"
	"fmt"

	"streaming/api/internal/shared"
	"streaming/api/internal/streaming/pb"
	"streaming/api/internal/streaming/response/model"

	"google.golang.org/grpc/metadata"
)

func (ctrl *Controller) GetResult(ctx context.Context, queryId int64) (model.ResponseDto, error) {
	ctx, span := ctrl.tracer.Start(ctx, "controller.GetResult")
	defer span.End()

	traceId := fmt.Sprintf("%s", span.SpanContext().TraceID())
	ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceId)

	respDb, err := ctrl.repo.FindOne(ctx, queryId)
	if err != nil {
		return model.ResponseDto{}, err
	}

	if respDb.Status == shared.ResponseStatusError || respDb.Status == shared.ResponseStatusFramerError || respDb.Status == shared.ResponseStatusDetectionError {
		return model.ResponseDto{
			QueryId: queryId,
			Status:  respDb.Status,
		}, nil
	}

	resp, err := ctrl.orchestrator.GetResult(ctx, &pb.GetResultReq{
		QueryId: queryId,
	})
	if err != nil {
		return model.ResponseDto{}, err
	}

	rawSources := resp.GetSources()
	sources := make([]model.Source, 0, len(rawSources))
	for _, source := range rawSources {
		sources = append(sources, model.Source{
			Id:   source.GetId(),
			Link: source.GetLink(),
		})
	}

	return model.ResponseDto{
		QueryId: queryId,
		Status:  shared.ResponseStatus(resp.GetStatus()),
		Sources: sources,
	}, nil
}
