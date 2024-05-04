package controller

import (
	"context"
	"strconv"
	"strings"

	"streaming/orchestrator/internal/streaming/model"
)

func (ctrl *Controller) GetResult(ctx context.Context, id int64) (model.ResponseDto, error) {
	ctx, err := getTraceCtx(ctx)
	if err != nil {
		return model.ResponseDto{}, err
	}

	ctx, span := ctrl.tracer.Start(ctx, "controller.GetResult")
	defer span.End()

	result, err := ctrl.repo.FindResult(ctx, id)
	if err != nil {
		return model.ResponseDto{}, err
	}

	respStatus, err := ctrl.repo.FindOne(ctx, id)
	if err != nil {
		return model.ResponseDto{}, err
	}

	resp := model.ResponseDto{
		Id:      id,
		Status:  respStatus,
		Sources: make([]model.Source, 0),
	}

	for _, detection := range result {
		fileIdString := strings.Split(strings.Split(detection.Filename, "_")[1], ".")[0]
		fileId, _ := strconv.ParseInt(fileIdString, 10, 64)
		resp.Sources = append(resp.Sources, model.Source{
			Id:   fileId,
			Link: detection.Link,
		})
	}

	return resp, nil
}
