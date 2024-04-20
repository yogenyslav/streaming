package kafka

import (
	"errors"

	"github.com/IBM/sarama"
	"github.com/yogenyslav/logger"
)

type AsyncProducer struct {
	brokers  []string
	producer sarama.AsyncProducer
}

func MustNewAsyncProducer(config *Config) *AsyncProducer {
	cfg := sarama.NewConfig()

	cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	cfg.Producer.RequiredAcks = sarama.WaitForAll

	cfg.Producer.Return.Successes = false
	cfg.Producer.Return.Errors = true

	asyncProducer, err := sarama.NewAsyncProducer(config.Brokers, cfg)

	if err != nil {
		logger.Panic(err)
	}

	go func() {
		for e := range asyncProducer.Errors() {
			logger.Error(e)
		}
	}()

	return &AsyncProducer{
		brokers:  config.Brokers,
		producer: asyncProducer,
	}
}

func (k *AsyncProducer) SendAsyncMessage(message *sarama.ProducerMessage) {
	k.producer.Input() <- message
}

func (k *AsyncProducer) Close() error {
	if err := k.producer.Close(); err != nil {
		return errors.Join(err, errors.New("kafka.Connector.Close"))
	}
	return nil
}
