package handler

import (
	"context"

	"streaming/orchestrator/internal/streaming/pb"
)

func (h *Handler) GetResult(ctx context.Context, params *pb.GetResultReq) (*pb.Response, error) {
	id := params.GetQueryId()
	resp, err := h.ctrl.GetResult(ctx, id)
	if err != nil {
		return nil, err
	}
	return convertToPb(resp), nil
}
