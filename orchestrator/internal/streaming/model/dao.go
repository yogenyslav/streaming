package model

import (
	"time"

	"streaming/orchestrator/internal/shared"
)

type ResponseDao struct {
	QueryId   int64                 `db:"query_id"`
	Status    shared.ResponseStatus `db:"status"`
	CreatedAt time.Time             `db:"created_at"`
	UpdatedAt time.Time             `db:"updated_at"`
}

type DetectionDao struct {
	QueryId  int64  `bson:"query_id"`
	Filename string `bson:"filename"`
	Lb       int    `bson:"lb"`
	Lt       int    `bson:"lt"`
	Rb       int    `bson:"rb"`
	Rt       int    `bson:"rt"`
	Link     string `bson:"link"`
}
