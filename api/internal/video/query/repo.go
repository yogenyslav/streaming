package query

import (
	"context"

	"streaming/internal/video/query/model"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yogenyslav/storage/postgres"
)

type Repo struct {
	pg *pgxpool.Pool
}

func NewRepo(pg *pgxpool.Pool) *Repo {
	return &Repo{
		pg: pg,
	}
}

const insertOne = `
	insert into query(type, source)
	values ($1, $2)
	returning id;
`

func (r *Repo) InsertOne(ctx context.Context, params model.Query) (int64, error) {
	return postgres.QueryPrimitive[int64](ctx, r.pg, insertOne, params.Type, params.Source)
}
