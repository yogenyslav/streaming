package response

import (
	"context"

	"streaming/internal/shared"
	"streaming/internal/video/response/model"

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
	insert into response(query_id, status)
	values ($1, $2);
`

func (r *Repo) InsertOne(ctx context.Context, params model.Response) error {
	_, err := r.pg.Exec(ctx, insertOne, params.QueryId, shared.StatusProcessing)
	return err
}

const updateOne = `
	update response
	set status = $1, updated_at = current_timestamp
	where query_id = $2;
`

func (r *Repo) UpdateOne(ctx context.Context, params model.Response) error {
	_, err := r.pg.Exec(ctx, updateOne, params.Status, params.QueryId)
	return err
}

const findOneById = `
	select query_id, status, created_at, updated_at
	from response
	where query_id = $1;
`

func (r *Repo) FindOneByQueryId(ctx context.Context, queryId int64) (model.Response, error) {
	return postgres.QueryStruct[model.Response](ctx, r.pg, findOneById, queryId)
}
