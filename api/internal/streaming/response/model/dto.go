package model

import (
	"streaming/api/internal/shared"
)

type Source struct {
	Id   int64  `json:"id"`
	Link string `json:"link"`
}

type ResponseDto struct {
	QueryId int64                 `json:"query_id"`
	Status  shared.ResponseStatus `json:"status"`
	Sources []Source              `json:"sources"`
}
