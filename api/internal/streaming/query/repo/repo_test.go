//go:build integration

package repo

import (
	"context"
	"testing"

	"streaming/api/config"
	"streaming/api/internal/streaming/query/model"
	"streaming/api/tests/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/postgres"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

var (
	cfg      = config.MustNew("../../../../config/config.test.yaml")
	exporter = tracetest.NewNoopExporter()
	provider = tracing.MustNewTraceProvider(exporter, "test-app")
	tracer   = otel.Tracer("test")
	db       = postgres.MustNew(cfg.Postgres, tracer)
)

func init() {
	otel.SetTracerProvider(provider)
}

func TestRepo_InsertOne(t *testing.T) {
	ctx := context.Background()
	defer database.TruncateTable(t, ctx, db.GetPool(), "query")
	repo := New(db)

	tests := []struct {
		name   string
		query  model.QueryDao
		wantId int64
	}{
		{
			name: "first insert, id must be 1",
			query: model.QueryDao{
				Source: "source_1",
			},
			wantId: 1,
		},
		{
			name: "second insert, id must be 2",
			query: model.QueryDao{
				Source: "source_2",
			},
			wantId: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, err := repo.InsertOne(ctx, tt.query)
			require.NoError(t, err)
			assert.Equal(t, tt.wantId, gotId)
		})
	}
}
