//go:build integration

//go:generate mockgen -source=../../pb/api_grpc.pb.go -destination=./mock/api.go -package=mock

package controller

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"streaming/api/config"
	"streaming/api/internal/shared"
	"streaming/api/internal/streaming/query/controller/mock"
	"streaming/api/internal/streaming/query/model"
	qr "streaming/api/internal/streaming/query/repo"
	rr "streaming/api/internal/streaming/response/repo"
	"streaming/api/tests/database"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/minios3"
	"github.com/yogenyslav/pkg/storage/postgres"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.uber.org/mock/gomock"
)

var (
	cfg      = config.MustNew("../../../../config/config.test.yaml")
	exporter = tracetest.NewNoopExporter()
	provider = tracing.MustNewTraceProvider(exporter, "test-app")
	tracer   = otel.Tracer("test")
	db       = postgres.MustNew(cfg.Postgres, tracer)
	s3       = minios3.MustNew(cfg.S3, tracer)
)

func init() {
	otel.SetTracerProvider(provider)
}

func TestController_InsertOne(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	orchestrator := mock.NewMockOrchestratorClient(ctrl)
	defer ctrl.Finish()

	videos := s3.ListObjects(ctx, shared.VideoBucket, minio.ListObjectsOptions{})
	videoCount := 0
	for range videos {
		videoCount++
	}

	defer database.TruncateTable(t, ctx, db.GetPool(), "query", "response")
	queryRepoTest := qr.New(db)
	responseRepoTest := rr.New(db)
	controller := New(queryRepoTest, responseRepoTest, orchestrator, s3, tracer)

	tests := []struct {
		name   string
		query  model.QueryCreate
		wantId int64
	}{
		{
			name: "insert 1, video",
			query: model.QueryCreate{
				Type:  shared.QueryTypeVideo,
				Video: bytes.NewReader([]byte("some data")),
				Size:  9,
				Name:  "test-video.mp4",
			},
			wantId: 1,
		},
		{
			name: "insert 2, video",
			query: model.QueryCreate{
				Type:  shared.QueryTypeVideo,
				Video: bytes.NewReader([]byte("more data")),
				Size:  9,
				Name:  "test-video2.mp4",
			},
			wantId: 2,
		},
		{
			name: "insert 3, stream",
			query: model.QueryCreate{
				Type:   shared.QueryTypeStream,
				Source: "rtsp://test.com",
			},
			wantId: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orchestrator.EXPECT().Process(gomock.Any(), gomock.Any())
			id, err := controller.InsertOne(ctx, tt.query)
			require.NoError(t, err)
			assert.Equal(t, tt.wantId, id)

			if tt.query.Type == shared.QueryTypeVideo {
				curVideos := s3.ListObjects(ctx, shared.VideoBucket, minio.ListObjectsOptions{})
				count := 0
				for range curVideos {
					count++
				}
				assert.Equal(t, count-videoCount, 1)

				split := strings.Split(tt.query.Name, ".")
				source := fmt.Sprintf("%s-%d.%s", tt.query.Name, tt.query.Ts, split[len(split)-1])
				err = s3.RemoveObject(ctx, shared.VideoBucket, source, minio.RemoveObjectOptions{})
				require.NoError(t, err)
			}
		})
	}
}
