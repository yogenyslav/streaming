package model

import (
	"time"

	"streaming/api/internal/shared"
)

type Response struct {
	QueryId   int64                 `db:"query_id"`
	Status    shared.ResponseStatus `db:"status"`
	CreatedAt time.Time             `db:"created_at"`
	UpdatedAt time.Time             `db:"updated_at"`
}

func (r Response) ToDto() ResponseDto {
	return ResponseDto{
		Status:    r.Status,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
