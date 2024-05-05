package model

import (
	"time"
)

type QueryDao struct {
	Id        int64     `db:"id"`
	Source    string    `db:"source"`
	CreatedAt time.Time `db:"created_at"`
}
