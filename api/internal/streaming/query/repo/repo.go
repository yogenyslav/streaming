package repo

import (
	"context"

	"streaming/api/internal/streaming/query/model"

	"github.com/yogenyslav/pkg/storage/postgres"
)

type Repo struct {
	db postgres.Database
}

func New(db postgres.Database) *Repo {
	return &Repo{
		db: db,
	}
}

const insertOne = `
	insert into query(source)
	values ($1)
	returning id;
`

func (r *Repo) InsertOne(ctx context.Context, params model.QueryDao) (int64, error) {
	var id int64
	err := r.db.Query(ctx, &id, insertOne, params.Source)
	return id, err
}
