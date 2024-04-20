package kafka

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/yogenyslav/logger"
)

type Consumer struct {
	brokers        []string
	SingleConsumer sarama.Consumer
}

func MustNewConsumer(config *Config) *Consumer {
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = false
	cfg.Consumer.Offsets.AutoCommit.Enable = true
	cfg.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer(config.Brokers, cfg)
	if err != nil {
		logger.Panic(err)
	}

	return &Consumer{
		brokers:        config.Brokers,
		SingleConsumer: consumer,
	}
}

func (consumer *Consumer) Subscribe(ctx context.Context, topic string, out chan<- *sarama.ConsumerMessage) error {
	partitions, err := consumer.SingleConsumer.Partitions(topic)
	if err != nil {
		logger.Errorf("error getting partitions: %v", err)
		return err
	}

	initialOffset := sarama.OffsetNewest

	for _, partition := range partitions {
		pc, err := consumer.SingleConsumer.ConsumePartition(topic, partition, initialOffset)
		if err != nil {
			logger.Errorf("error consuming partition: %v", err)
			return err
		}

		go consume(ctx, pc, partition, out)
	}
	return nil
}

func consume(ctx context.Context, pc sarama.PartitionConsumer, partition int32, out chan<- *sarama.ConsumerMessage) {
	logger.Infof("consumer started for partition %d", partition)
	for {
		select {
		case <-ctx.Done():
			pc.Close()
			logger.Infof("consumer closed for partition %d", partition)
			return
		case message := <-pc.Messages():
			out <- message
			logger.Infof(
				"[kafka] message claimed, topic=%s, partition=%d, key=%s offset=%d",
				message.Topic, message.Partition, message.Key, message.Offset,
			)
		}
	}
}
