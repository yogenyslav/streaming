package repo

import (
	"context"

	"streaming/orchestrator/internal/shared"
	"streaming/orchestrator/internal/streaming/model"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/storage/postgres"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repo struct {
	db    postgres.Database
	mongo *mongo.Client
}

func New(db postgres.Database, mongo *mongo.Client) *Repo {
	return &Repo{
		db:    db,
		mongo: mongo,
	}
}

const updateOne = `
	update response 
	set status = $2, updated_at=current_timestamp
	where query_id = $1;
`

func (r *Repo) UpdateOne(ctx context.Context, queryId int64, status shared.ResponseStatus) error {
	_, err := r.db.Exec(ctx, updateOne, queryId, status)
	return err
}

const findOne = `
	select status
	from response
	where query_id = $1;
`

func (r *Repo) FindOne(ctx context.Context, queryId int64) (shared.ResponseStatus, error) {
	var status shared.ResponseStatus
	err := r.db.Query(ctx, &status, findOne, queryId)
	return status, err
}

func (r *Repo) FindResult(ctx context.Context, queryId int64) ([]model.DetectionDao, error) {
	var res []model.DetectionDao
	col := r.mongo.Database("dev").Collection("frames")

	cursor, err := col.Find(ctx, bson.M{"query_id": queryId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &res); err != nil {
		return nil, err
	}

	log.Info().Int64("queryId", queryId).Int("count", len(res)).Msg("find result")

	return res, nil
}
