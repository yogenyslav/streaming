package repo

import (
	"context"
	"testing"

	"streaming/api/config"
	"streaming/api/internal/shared"
	"streaming/api/internal/streaming/response/model"
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
	defer database.TruncateTable(t, ctx, db.GetPool(), "response")
	repo := New(db)

	tests := []struct {
		name       string
		response   model.ResponseDao
		wantErr    bool
		wantId     int64
		wantStatus shared.ResponseStatus
	}{
		{
			name: "success, insert 1",
			response: model.ResponseDao{
				QueryId: 1,
			},
			wantId:     1,
			wantStatus: shared.ResponseStatusInit,
		},
		{
			name: "success, insert 2",
			response: model.ResponseDao{
				QueryId: 2,
			},
			wantId:     2,
			wantStatus: shared.ResponseStatusInit,
		},
		{
			name: "fail, insert same query id",
			response: model.ResponseDao{
				QueryId: 2,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.InsertOne(ctx, tt.response)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			resp := model.ResponseDao{}
			err = db.Query(ctx, &resp, "select * from response where query_id = $1", tt.response.QueryId)
			require.NoError(t, err)
			assert.Equal(t, tt.wantId, resp.QueryId)
			assert.Equal(t, tt.wantStatus, resp.Status)
		})
	}
}

func TestRepo_FindOne(t *testing.T) {
	ctx := context.Background()
	defer database.TruncateTable(t, ctx, db.GetPool(), "response")
	repo := New(db)

	_, err := db.Exec(ctx, "insert into response (query_id) values ($1)", 1)
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name: "success",
			id:   1,
		},
		{
			name:    "fail, not found",
			id:      2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := repo.FindOne(ctx, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.id, resp.QueryId)
			assert.Equal(t, shared.ResponseStatusInit, resp.Status)
		})
	}
}

func TestRepo_UpdateOne(t *testing.T) {
	ctx := context.Background()
	defer database.TruncateTable(t, ctx, db.GetPool(), "response")
	repo := New(db)

	_, err := db.Exec(ctx, "insert into response (query_id) values ($1)", 1)
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      int64
		status  shared.ResponseStatus
		wantErr bool
	}{
		{
			name:   "success, update to framer processing",
			id:     1,
			status: shared.ResponseStatusFramerProcessing,
		},
		{
			name:   "success, update to error",
			id:     1,
			status: shared.ResponseStatusError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateOne(ctx, tt.id, tt.status)
			assert.NoError(t, err)

			resp := model.ResponseDao{}
			err = db.Query(ctx, &resp, "select * from response where query_id = $1", tt.id)
			require.NoError(t, err)
			assert.Equal(t, tt.status, resp.Status)
		})
	}
}
