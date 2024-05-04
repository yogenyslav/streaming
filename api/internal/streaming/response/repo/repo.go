package repo

import (
	"context"

	"streaming/api/internal/shared"
	"streaming/api/internal/streaming/response/model"

	"github.com/yogenyslav/pkg/storage/postgres"
)

type Repo struct {
	db postgres.Database
}

func New(db postgres.Database) *Repo {
	return &Repo{db: db}
}

const insertOne = `
	insert into response (query_id)
	values ($1);
`

func (r *Repo) InsertOne(ctx context.Context, params model.ResponseDao) error {
	_, err := r.db.Exec(ctx, insertOne, params.QueryId)
	return err
}

const findOne = `
	select query_id, status
	from response
	where query_id=$1;
`

func (r *Repo) FindOne(ctx context.Context, id int64) (model.ResponseDao, error) {
	var resp model.ResponseDao
	err := r.db.Query(ctx, &resp, findOne, id)
	return resp, err
}

const updateOne = `
	update response 
	set status=$2, updated_at=current_timestamp
	where query_id=$1;
`

func (r *Repo) UpdateOne(ctx context.Context, id int64, status shared.ResponseStatus) error {
	_, err := r.db.Exec(ctx, updateOne, id, status)
	return err
}
