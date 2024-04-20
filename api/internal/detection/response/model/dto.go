package model

import (
	"time"

	"streaming/api/internal/shared"
)

type ResponseCreateReq struct {
	QueryId int64
}

type ResponseUpdateReq struct {
	QueryId int64
	Status  shared.ResponseStatus
}

type ResponseDto struct {
	Status    shared.ResponseStatus `json:"status"`
	Sources   []string              `json:"sources,omitempty"`
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
}
