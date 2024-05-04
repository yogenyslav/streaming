package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"streaming/orchestrator/config"
	srvconf "streaming/orchestrator/internal/server/config"
	"streaming/orchestrator/internal/streaming/controller"
	"streaming/orchestrator/internal/streaming/handler"
	"streaming/orchestrator/internal/streaming/pb"
	"streaming/orchestrator/internal/streaming/repo"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/infrastructure/kafka"
	"github.com/yogenyslav/pkg/infrastructure/prom"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/minios3"
	"github.com/yogenyslav/pkg/storage/postgres"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type Server struct {
	cfg    *config.Config `yaml:"server"`
	conn   *grpc.Server
	db     postgres.Postgres
	s3     minios3.S3
	mongo  *mongo.Client
	tracer trace.Tracer
}

func New(cfg *config.Config) *Server {
	var grpcOpts []grpc.ServerOption
	conn := grpc.NewServer(grpcOpts...)

	tracing.MustSetupOTel(fmt.Sprintf("%s:%d", cfg.Tracing.Host, cfg.Tracing.Port), "orchestrator")
	tracer := otel.Tracer("orchestrator")

	clusterAdmin, err := sarama.NewClusterAdmin([]string{"kafka:29091"}, sarama.NewConfig())
	if err != nil {
		panic(err)
	}
	for _, topic := range cfg.Kafka.Topics {
		err = clusterAdmin.CreateTopic(topic.Name, &sarama.TopicDetail{
			NumPartitions:     topic.Partitions,
			ReplicationFactor: 1,
		}, false)
		if err != nil && !errors.Is(err, sarama.ErrTopicAlreadyExists) {
			panic(err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d/%s", cfg.Mongo.Host, cfg.Mongo.Port, cfg.Mongo.Db)))
	if err != nil {
		panic(err)
	}

	return &Server{
		cfg:    cfg,
		conn:   conn,
		db:     postgres.MustNew(cfg.Postgres, tracer),
		s3:     minios3.MustNew(cfg.S3, tracer),
		mongo:  mongoClient,
		tracer: tracer,
	}
}

func (s *Server) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer s.mongo.Disconnect(ctx)

	if err := s.s3.CreateBuckets(ctx); err != nil {
		panic(err)
	}

	consumer := kafka.MustNewConsumer(s.cfg.Kafka, true, 5)
	defer consumer.SingleConsumer.Close()

	producer := kafka.MustNewAsyncProducer(s.cfg.Kafka, sarama.NewRandomPartitioner, sarama.WaitForAll)
	defer producer.Close()

	streamingRepo := repo.New(s.db, s.mongo)
	streamingController := controller.New(ctx, streamingRepo, s.s3, consumer, producer, s.tracer)
	streamingHandler := handler.New(streamingController)
	pb.RegisterOrchestratorServer(s.conn, streamingHandler)

	go s.listenGrpc(s.cfg.Server)
	go prom.HandlePrometheus(s.cfg.Prometheus.Host, s.cfg.Prometheus.Port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	cancel()
	s.conn.GracefulStop()
	log.Info().Msg("server was gracefully stopped")
	os.Exit(0)
}

func (s *Server) listenGrpc(cfg *srvconf.ServerConfig) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		panic(err)
	}

	if err = s.conn.Serve(lis); err != nil {
		log.Error().Err(err).Msg("error while listening, stopping grpc server")
	}
}
