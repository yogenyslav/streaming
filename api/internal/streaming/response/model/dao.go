package model

import (
	"time"

	"streaming/api/internal/shared"
)

type ResponseDao struct {
	QueryId   int64                 `db:"query_id"`
	Status    shared.ResponseStatus `db:"status"`
	CreatedAt time.Time             `db:"created_at"`
	UpdatedAt time.Time             `db:"updated_at"`
}
