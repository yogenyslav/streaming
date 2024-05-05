package controller

import (
	"context"

	"streaming/orchestrator/internal/shared"

	"github.com/rs/zerolog/log"
)

func (ctrl *Controller) CancelById(ctx context.Context, id int64) (bool, error) {
	ctx, err := getTraceCtx(ctx)
	if err != nil {
		return false, err
	}

	ctx, span := ctrl.tracer.Start(ctx, "controller.CancelById")
	defer span.End()

	log.Info().Int64("queryId", id).Msg("was canceled")
	if err = ctrl.repo.UpdateOne(ctx, id, shared.ResponseStatusCanceled); err != nil {
		log.Error().Err(err).Msg("update response error")
		return false, err
	}

	ctrl.mu.Lock()
	cancel, ok := ctrl.processing[id]
	if !ok {
		ctrl.mu.Unlock()
		return false, nil
	}
	cancel()
	delete(ctrl.processing, id)
	ctrl.mu.Unlock()
	return true, nil
}
