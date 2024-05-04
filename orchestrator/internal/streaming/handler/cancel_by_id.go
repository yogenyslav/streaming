package handler

import (
	"context"

	"streaming/orchestrator/internal/streaming/pb"
)

func (h *Handler) Cancel(ctx context.Context, params *pb.CancelReq) (*pb.CancelResp, error) {
	id := params.GetQueryId()
	success, err := h.ctrl.CancelById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pb.CancelResp{Success: success}, nil
}
