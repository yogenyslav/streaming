package server

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"streaming/orchestrator/config"
	"streaming/orchestrator/internal/orchestrator"
	"streaming/orchestrator/internal/pb"
	"streaming/orchestrator/pkg/infrastructure/kafka"

	"github.com/IBM/sarama"
	"github.com/yogenyslav/logger"
	"google.golang.org/grpc"
)

type Server struct {
	conn     *grpc.Server
	cfg      *config.Config
	consumer *kafka.Consumer
	producer *kafka.AsyncProducer
}

func New(cfg *config.Config) *Server {
	var grpcOpts []grpc.ServerOption
	conn := grpc.NewServer(grpcOpts...)

	clusterAdmin, err := sarama.NewClusterAdmin([]string{"kafka1:29091"}, sarama.NewConfig())
	if err != nil {
		logger.Panic(err)
	}
	for _, topic := range cfg.Kafka.Topics {
		err = clusterAdmin.CreateTopic(topic.Name, &sarama.TopicDetail{
			NumPartitions:     topic.Partitions,
			ReplicationFactor: 1,
		}, false)
		if err != nil && !errors.Is(err, sarama.ErrTopicAlreadyExists) {
			logger.Panic(err)
		}
	}

	return &Server{
		conn:     conn,
		cfg:      cfg,
		consumer: kafka.MustNewConsumer(cfg.Kafka),
		producer: kafka.MustNewAsyncProducer(cfg.Kafka),
	}
}

func (s *Server) Run() {
	controller := orchestrator.NewController(s.consumer, s.producer)
	handler := orchestrator.NewHandler(controller)

	pb.RegisterOrchestratorServer(s.conn, handler)

	go s.listen()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	if err := s.consumer.SingleConsumer.Close(); err != nil {
		logger.Warnf("failed to close kafka consumer: %v", err)
	}
	if err := s.producer.Close(); err != nil {
		logger.Warnf("failed to close kafka producer: %v", err)
	}
	s.conn.GracefulStop()

	logger.Info("graceful shutdown finished")
	os.Exit(0)
}

func (s *Server) listen() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", s.cfg.Server.Port))
	if err != nil {
		logger.Panic(err)
	}

	if err = s.conn.Serve(lis); err != nil {
		logger.Panicf("unexpected error, stopping grpc server: %v", err)
	}
}
