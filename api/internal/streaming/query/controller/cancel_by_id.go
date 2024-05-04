package controller

import (
	"context"
	"fmt"

	"streaming/api/internal/shared"
	"streaming/api/internal/streaming/pb"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
)

func (ctrl *Controller) CancelById(ctx context.Context, id int64) error {
	ctx, span := ctrl.tracer.Start(ctx, "controller.CancelById")
	defer span.End()

	traceId := fmt.Sprintf("%s", span.SpanContext().TraceID())
	ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceId)

	resp, err := ctrl.orchestrator.Cancel(ctx, &pb.CancelReq{
		QueryId: id,
	})
	if err != nil || !resp.Success {
		log.Error().Err(err).Msg("failed to cancel query")
		return shared.ErrCancelQuery
	}

	return nil
}
