package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"streaming/orchestrator/internal/shared"
	"streaming/orchestrator/internal/streaming/model"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/loctime"
)

func (ctrl *Controller) Process(ctx context.Context, params model.QueryDto) (model.ResponseDto, error) {
	var ok bool

	ctx, err := getTraceCtx(ctx)
	if err != nil {
		return model.ResponseDto{}, err
	}

	ctx, span := ctrl.tracer.Start(ctx, "controller.Process")
	defer span.End()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ctrl.mu.Lock()
	ctrl.processing[params.Id] = cancel
	ctrl.mu.Unlock()

	id, err := ctrl.processFramer(ctx, params)
	if err != nil {
		ctrl.mu.Lock()
		cancel, ok = ctrl.processing[id]
		if ok {
			cancel()
		}
		delete(ctrl.processing, id)
		ctrl.mu.Unlock()

		return model.ResponseDto{}, err
	}

	ctx, cancel = context.WithCancel(ctx)
	defer cancel()
	ctrl.mu.Lock()
	ctrl.processing[id] = cancel
	ctrl.mu.Unlock()

	id, err = ctrl.processDetection(ctx, id)
	if err != nil {
		ctrl.mu.Lock()
		cancel, ok = ctrl.processing[id]
		if ok {
			cancel()
		}
		delete(ctrl.processing, id)
		ctrl.mu.Unlock()

		return model.ResponseDto{}, err
	}

	ctx, cancel = context.WithCancel(ctx)
	defer cancel()
	ctrl.mu.Lock()
	ctrl.processing[id] = cancel
	ctrl.mu.Unlock()

	id, err = ctrl.processResponser(ctx, id)
	if err != nil {
		ctrl.mu.Lock()
		cancel, ok = ctrl.processing[id]
		if ok {
			cancel()
		}
		delete(ctrl.processing, id)
		ctrl.mu.Unlock()

		return model.ResponseDto{}, err
	}

	ctrl.mu.Lock()
	cancel, ok = ctrl.processing[id]
	if ok {
		cancel()
	}
	delete(ctrl.processing, id)
	ctrl.mu.Unlock()

	return model.ResponseDto{
		Id:     id,
		Status: shared.ResponseStatusSuccess,
	}, nil
}

func (ctrl *Controller) processFramer(ctx context.Context, params model.QueryDto) (int64, error) {
	ctx, span := ctrl.tracer.Start(ctx, "controller.processFramer")
	defer span.End()

	data, err := json.Marshal(params)
	if err != nil {
		log.Error().Err(err).Msg("processFramer marshal error")
		return params.Id, err
	}

	ctrl.producer.SendAsyncMessage(ctx, &sarama.ProducerMessage{
		Topic:     shared.FramerTopic,
		Key:       sarama.StringEncoder(fmt.Sprintf("%d", params.Id)),
		Value:     sarama.ByteEncoder(data),
		Timestamp: loctime.GetLocalTime(),
	})

	if err = ctrl.repo.UpdateOne(ctx, params.Id, shared.ResponseStatusFramerProcessing); err != nil {
		log.Error().Err(err).Msg("processFramer update response error")
		return params.Id, err
	}

	for {
		select {
		case <-ctx.Done():
			return params.Id, nil
		case rawMessage := <-ctrl.framerRes:
			message := model.ResultMessage{}
			if err = json.Unmarshal(rawMessage.Value, &message); err != nil {
				log.Panic().Err(err)
			}

			if message.Message == shared.ErrorMessage {
				if err = ctrl.repo.UpdateOne(ctx, params.Id, shared.ResponseStatusFramerError); err != nil {
					log.Error().Err(err).Msg("processFramer update response error")
					return message.Id, err
				}
				return message.Id, shared.ErrFramer
			}

			return message.Id, nil
		}
	}
}

func (ctrl *Controller) processDetection(ctx context.Context, id int64) (int64, error) {
	ctx, span := ctrl.tracer.Start(ctx, "controller.processDetection")
	defer span.End()

	if err := ctrl.repo.UpdateOne(ctx, id, shared.ResponseStatusDetectionProcessing); err != nil {
		log.Error().Err(err).Msg("processDetection update response error")
		return id, err
	}

	for {
		select {
		case <-ctx.Done():
			req := struct {
				QueryId int64 `json:"query_id"`
				Cancel  bool  `json:"cancel"`
			}{
				QueryId: id,
				Cancel:  true,
			}
			data, _ := json.Marshal(req)
			ctrl.producer.SendAsyncMessage(ctx, &sarama.ProducerMessage{
				Topic:     shared.DetectionCancelTopic,
				Key:       sarama.StringEncoder(fmt.Sprintf("%d", id)),
				Value:     sarama.ByteEncoder(data),
				Timestamp: loctime.GetLocalTime(),
			})

			return id, nil
		case rawMessage := <-ctrl.detectionRes:
			var err error
			message := model.ResultMessage{}
			if err = json.Unmarshal(rawMessage.Value, &message); err != nil {
				log.Panic().Err(err)
			}

			if message.Message == shared.ErrorMessage {
				if err = ctrl.repo.UpdateOne(ctx, id, shared.ResponseStatusDetectionError); err != nil {
					log.Error().Err(err).Msg("processDetection update response error")
					return message.Id, err
				}
				return message.Id, shared.ErrDetection
			}

			return message.Id, nil
		}
	}
}

func (ctrl *Controller) processResponser(ctx context.Context, id int64) (int64, error) {
	ctx, span := ctrl.tracer.Start(ctx, "controller.processResponser")
	defer span.End()

	if err := ctrl.repo.UpdateOne(ctx, id, shared.ResponseStatusResponserProcessing); err != nil {
		log.Error().Err(err).Msg("processResponser update response error")
		return id, err
	}

	for {
		select {
		case <-ctx.Done():
			return id, nil
		case rawMessage := <-ctrl.responserRes:
			var err error
			message := model.ResultMessage{}
			if err = json.Unmarshal(rawMessage.Value, &message); err != nil {
				log.Panic().Err(err)
			}

			return message.Id, nil
		}
	}
}
