package handler

import (
	"context"

	"streaming/orchestrator/internal/shared"
	"streaming/orchestrator/internal/streaming/model"
	"streaming/orchestrator/internal/streaming/pb"
)

func (h *Handler) Process(ctx context.Context, params *pb.Query) (*pb.Response, error) {
	query := model.QueryDto{
		Id:     params.GetId(),
		Source: params.GetSource(),
		Type:   shared.QueryType(params.GetType()),
	}

	resp, err := h.ctrl.Process(ctx, query)
	if err != nil {
		return &pb.Response{
			Id:     params.GetId(),
			Status: pb.ResponseStatus(shared.ResponseStatusError),
		}, err
	}

	return convertToPb(resp), nil
}
