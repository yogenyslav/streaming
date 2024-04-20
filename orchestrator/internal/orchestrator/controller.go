package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"streaming/orchestrator/internal/orchestrator/model"
	"streaming/orchestrator/internal/pb"
	"streaming/orchestrator/internal/shared"
	"streaming/orchestrator/pkg"
	"streaming/orchestrator/pkg/infrastructure/kafka"

	"github.com/IBM/sarama"
	"github.com/yogenyslav/logger"
)

type Query struct {
	Cancel context.CancelFunc
	Status pb.ResponseStatus
}

type Controller struct {
	consumer   *kafka.Consumer
	producer   *kafka.AsyncProducer
	processing map[int64]Query
	mu         sync.Mutex
	framer     chan *sarama.ConsumerMessage
	detection  chan *sarama.ConsumerMessage
}

func NewController(consumer *kafka.Consumer, producer *kafka.AsyncProducer) *Controller {
	return &Controller{
		consumer:  consumer,
		producer:  producer,
		framer:    make(chan *sarama.ConsumerMessage),
		detection: make(chan *sarama.ConsumerMessage),
	}
}

func (ctrl *Controller) Process(ctx context.Context, params *pb.Query) (pb.ResponseStatus, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer ctrl.cancel(params.Id)

	status := pb.ResponseStatus_INIT_PROCESSING

	ctrl.mu.Lock()
	ctrl.processing[params.GetId()] = Query{
		Cancel: cancel,
		Status: status,
	}
	ctrl.mu.Unlock()

	if err := ctrl.consumer.Subscribe(ctx, shared.TopicFramerResult, ctrl.framer); err != nil {
		logger.Errorf("failed to subcribe to framer: %v", err)
		return pb.ResponseStatus_ERROR, err
	}

	if err := ctrl.consumer.Subscribe(ctx, shared.TopicDetectionResult, ctrl.detection); err != nil {
		logger.Errorf("failed to subcribe to detection: %v", err)
		return pb.ResponseStatus_ERROR, err
	}

	data, err := json.Marshal(params)
	if err != nil {
		logger.Error(err)
		return pb.ResponseStatus_ERROR, err
	}

	ctrl.producer.SendAsyncMessage(&sarama.ProducerMessage{
		Topic:     shared.TopicFramer,
		Key:       sarama.StringEncoder(strconv.FormatInt(params.Id, 10)),
		Value:     sarama.ByteEncoder(data),
		Timestamp: pkg.GetLocalTime(),
	})

	status = pb.ResponseStatus_FRAMER_PROCESSING
	ctrl.updateStatus(params.GetId(), status)

	for {
		select {
		case <-ctx.Done():
			logger.Info("finished processing")
			return status, nil
		case message := <-ctrl.framer:
			frame := model.RawFrame{}
			if err = json.Unmarshal(message.Value, &frame); err != nil {
				status = pb.ResponseStatus_FRAMER_ERROR
				logger.Errorf("failed to unmarshal framer result: %v", err)
				return status, err
			}

			if frame.QueryId != params.Id {
				continue
			}

			ctrl.producer.SendAsyncMessage(&sarama.ProducerMessage{
				Topic:     shared.TopicDetection,
				Key:       sarama.StringEncoder(fmt.Sprintf("%d-%d", frame.QueryId, frame.Id)),
				Value:     sarama.ByteEncoder(message.Value),
				Timestamp: pkg.GetLocalTime(),
			})

			if frame.IsLast {
				status = pb.ResponseStatus_FRAMER_SUCCESS
				ctrl.updateStatus(params.GetId(), status)
			}
		case message := <-ctrl.detection:
			frame := model.ProcessedFrame{}
			if err = json.Unmarshal(message.Value, &frame); err != nil {
				status = pb.ResponseStatus_DETECTION_ERROR
				logger.Errorf("failed to unmarshal detection result: %v", err)
				return status, err
			}

			if frame.QueryId != params.Id {
				continue
			}

			ctrl.producer.SendAsyncMessage(&sarama.ProducerMessage{
				Topic:     shared.TopicResponser,
				Key:       sarama.StringEncoder(fmt.Sprintf("%d-%d", frame.QueryId, frame.Id)),
				Value:     sarama.ByteEncoder(message.Value),
				Timestamp: pkg.GetLocalTime(),
			})

			if frame.IsLast {
				status = pb.ResponseStatus_DETECTION_SUCCESS
				ctrl.updateStatus(params.GetId(), status)
			}
		}
	}
}

func (ctrl *Controller) cancel(id int64) {
	ctrl.mu.Lock()
	defer ctrl.mu.Unlock()

	query, ok := ctrl.processing[id]
	if !ok {
		logger.Errorf("query %d was not processing", id)
		return
	}

	query.Cancel()
	delete(ctrl.processing, id)
}

func (ctrl *Controller) updateStatus(id int64, status pb.ResponseStatus) {
	status = pb.ResponseStatus_FRAMER_SUCCESS
	ctrl.mu.Lock()
	query := ctrl.processing[id]
	query.Status = status
	ctrl.processing[id] = query
	ctrl.mu.Unlock()
}
