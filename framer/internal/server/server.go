package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"streaming/framer/config"
	"streaming/framer/internal/model"
	"streaming/framer/internal/pb"
	"streaming/framer/internal/shared"
	"streaming/framer/pkg"
	"streaming/framer/pkg/infrastructure/kafka"

	"github.com/IBM/sarama"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/yogenyslav/logger"
)

type Server struct {
	cfg      *config.Config
	consumer *kafka.Consumer
	producer *kafka.AsyncProducer
}

func New(cfg *config.Config) *Server {
	return &Server{
		cfg:      cfg,
		consumer: kafka.MustNewConsumer(cfg.Kafka),
		producer: kafka.MustNewAsyncProducer(cfg.Kafka),
	}
}

func (s *Server) Run() {
	defer s.producer.Close()
	defer s.consumer.SingleConsumer.Close()

	out := make(chan *sarama.ConsumerMessage)
	err := s.consumer.Subscribe(context.Background(), shared.TopicFramer, out)
	if err != nil {
		logger.Panic(err)
	}

	for {
		select {
		case message := <-out:
			go func() {
				query := pb.Query{}
				if err := json.Unmarshal(message.Value, &query); err != nil {
					logger.Errorf("failed to unmarshal query: %v", err)
					frame := model.RawFrame{
						Error: err,
					}
					data, _ := json.Marshal(frame)
					s.producer.SendAsyncMessage(&sarama.ProducerMessage{
						Topic:     shared.TopicFramerResult,
						Key:       sarama.StringEncoder(strconv.FormatInt(query.GetId(), 10)),
						Value:     sarama.ByteEncoder(data),
						Timestamp: pkg.GetLocalTime(),
					})
				}

				info, err := ffmpeg.Probe(query.GetSource())
				buf := bytes.NewBuffer(nil)
				err = ffmpeg.Input(query.GetSource()).
					Output("pipe:", ffmpeg.KwArgs{"format": "image2", "vcodec": "mjpeg"}).
					WithOutput(buf, os.Stdout).
					Run()
				if err != nil {
					logger.Errorf("failed to read frames: %v", err)
					frame := model.RawFrame{
						Error: err,
					}
					data, _ := json.Marshal(frame)
					s.producer.SendAsyncMessage(&sarama.ProducerMessage{
						Topic:     shared.TopicFramerResult,
						Key:       sarama.StringEncoder(strconv.FormatInt(query.GetId(), 10)),
						Value:     sarama.ByteEncoder(data),
						Timestamp: pkg.GetLocalTime(),
					})
				}

			}()
		}
	}
}

func getVideoSize(fileName string) (int, int, error) {
	data, err := ffmpeg.Probe(fileName)
	if err != nil {
		logger.Error(err)
		return 0, 0, err
	}
	logger.Infof("got video info: %s", data)
	type VideoInfo struct {
		Streams []struct {
			CodecType string `json:"codec_type"`
			Width     int
			Height    int
		} `json:"streams"`
	}
	vInfo := &VideoInfo{}
	err = json.Unmarshal([]byte(data), vInfo)
	if err != nil {
		logger.Error(err)
		return 0, 0, err
	}
	for _, s := range vInfo.Streams {
		if s.CodecType == "video" {
			return s.Width, s.Height, nil
		}
	}
	return 0, 0, errors.New("unable to specify size")
}
