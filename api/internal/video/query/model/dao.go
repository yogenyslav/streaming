package model

import (
	"time"

	"streaming/internal/shared"
)

type Query struct {
	Id        int64            `db:"id"`
	Type      shared.QueryType `db:"type"`
	Source    string           `db:"source"`
	CreatedAt time.Time        `db:"created_at"`
}
