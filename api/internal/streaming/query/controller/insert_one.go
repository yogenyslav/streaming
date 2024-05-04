package controller

import (
	"context"
	"fmt"
	"strings"

	"streaming/api/internal/shared"
	"streaming/api/internal/streaming/pb"
	"streaming/api/internal/streaming/query/model"
	respmodel "streaming/api/internal/streaming/response/model"

	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (ctrl *Controller) InsertOne(ctx context.Context, params model.QueryCreate) (int64, error) {
	ctx, span := ctrl.tracer.Start(ctx, "Controller.InsertOne")
	defer span.End()

	query := model.QueryDao{
		Source: params.Source,
	}

	if params.Type == shared.QueryTypeVideo {
		split := strings.Split(params.Name, ".")
		source := fmt.Sprintf("%s-%d.%s", params.Name, params.Ts, split[len(split)-1])
		_, err := ctrl.s3.PutObject(ctx, shared.VideoBucket, source, params.Video, params.Size, minio.PutObjectOptions{})
		if err != nil {
			return 0, err
		}
		query.Source = source
	}

	id, err := ctrl.qr.InsertOne(ctx, query)
	if err != nil {
		return 0, err
	}

	resp := respmodel.ResponseDao{QueryId: id}
	if err = ctrl.rr.InsertOne(ctx, resp); err != nil {
		return 0, err
	}

	query.Id = id
	go ctrl.process(query, params.Type)

	return id, nil
}

func (ctrl *Controller) process(query model.QueryDao, t shared.QueryType) {
	ctx, span := ctrl.tracer.Start(context.Background(), "controller.process")
	defer span.End()

	traceId := fmt.Sprintf("%s", span.SpanContext().TraceID())
	ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceId)

	resp, err := ctrl.orchestrator.Process(ctx, &pb.Query{
		Id:     query.Id,
		Source: query.Source,
		Type:   pb.QueryType(t),
	})
	if err != nil {
		grpcErr := status.Convert(err)
		if grpcErr.Code() != codes.Canceled {
			log.Error().Err(err).Msg("failed to process query")
		} else {
			log.Info().Msg("query processing cancelled")
			return
		}
	}

	log.Info().Int64("id", query.Id).Any("status", resp.GetStatus()).Msg("query processed")

	if err = ctrl.rr.UpdateOne(ctx, query.Id, shared.ResponseStatus(resp.GetStatus())); err != nil {
		log.Error().Err(err).Msg("failed to update response status")
		return
	}

	// TODO: save resp to cache
}
