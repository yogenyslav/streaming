package model

import (
	"streaming/orchestrator/internal/shared"
)

type QueryDto struct {
	Id     int64            `json:"id"`
	Source string           `json:"source"`
	Type   shared.QueryType `json:"type"`
	Cancel bool             `json:"cancel"`
}

type Source struct {
	Id   int64  `json:"id"`
	Link string `json:"link"`
}

type ResponseDto struct {
	Id      int64                 `json:"id"`
	Status  shared.ResponseStatus `json:"status"`
	Sources []Source              `json:"sources"`
}

type ResultMessage struct {
	Id      int64  `json:"id"`
	Message string `json:"message"`
}
