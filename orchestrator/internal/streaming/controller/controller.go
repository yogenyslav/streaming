package controller

import (
	"context"
	"sync"

	"streaming/orchestrator/internal/shared"
	"streaming/orchestrator/internal/streaming/model"

	"github.com/IBM/sarama"
	"github.com/yogenyslav/pkg/infrastructure/kafka"
	"github.com/yogenyslav/pkg/storage/minios3"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

type streamingRepo interface {
	UpdateOne(ctx context.Context, queryId int64, status shared.ResponseStatus) error
	FindOne(ctx context.Context, queryId int64) (shared.ResponseStatus, error)
	FindResult(ctx context.Context, queryId int64) ([]model.DetectionDao, error)
}

type Controller struct {
	repo         streamingRepo
	s3           minios3.S3
	consumer     *kafka.Consumer
	producer     *kafka.AsyncProducer
	tracer       trace.Tracer
	framerRes    chan *sarama.ConsumerMessage
	detectionRes chan *sarama.ConsumerMessage
	responserRes chan *sarama.ConsumerMessage
	mu           sync.Mutex
	processing   map[int64]context.CancelFunc
}

func New(ctx context.Context, repo streamingRepo, s3 minios3.S3, consumer *kafka.Consumer, producer *kafka.AsyncProducer, tracer trace.Tracer) *Controller {
	framerRes := make(chan *sarama.ConsumerMessage)
	err := consumer.Subscribe(ctx, shared.FramerResTopic, framerRes)
	if err != nil {
		panic(err)
	}

	detectionRes := make(chan *sarama.ConsumerMessage)
	err = consumer.Subscribe(ctx, shared.DetectionResTopic, detectionRes)
	if err != nil {
		panic(err)
	}

	responserRes := make(chan *sarama.ConsumerMessage)
	err = consumer.Subscribe(ctx, shared.ResponserResTopic, responserRes)
	if err != nil {
		panic(err)
	}

	return &Controller{
		repo:         repo,
		s3:           s3,
		consumer:     consumer,
		producer:     producer,
		tracer:       tracer,
		framerRes:    framerRes,
		detectionRes: detectionRes,
		responserRes: responserRes,
		processing:   make(map[int64]context.CancelFunc),
	}
}

func getTraceCtx(ctx context.Context) (context.Context, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	traceIdString := md["x-trace-id"][0]
	traceId, err := trace.TraceIDFromHex(traceIdString)
	if err != nil {
		return ctx, err
	}

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceId,
	})
	ctx = trace.ContextWithSpanContext(ctx, spanCtx)
	return ctx, nil
}
