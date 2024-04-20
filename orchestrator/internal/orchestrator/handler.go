package orchestrator

import (
	"context"

	"streaming/orchestrator/internal/pb"
)

type orchestratorController interface {
	Process(ctx context.Context, params *pb.Query) (pb.ResponseStatus, error)
}

type Handler struct {
	pb.UnimplementedOrchestratorServer
	controller orchestratorController
}

func NewHandler(controller orchestratorController) *Handler {
	return &Handler{
		controller: controller,
	}
}

func (h *Handler) Process(ctx context.Context, params *pb.Query) (*pb.Response, error) {
	status, err := h.controller.Process(ctx, params)
	return &pb.Response{
		Status: status,
	}, err
}

func (h *Handler) GetProcessed(ctx context.Context, params *pb.ProcessedReq) (*pb.ProcessedResp, error) {
	return nil, nil
}

func (h *Handler) Cancel(ctx context.Context, params *pb.CancelReq) (*pb.CancelResp, error) {
	return nil, nil
}
