package model

import (
	"io"

	"streaming/api/internal/shared"
)

type QueryCreate struct {
	Type   shared.QueryType
	Source string
	Name   string
	Video  io.Reader
	Size   int64
	Ts     int64
}

type QueryStreamReq struct {
	Source string `json:"source"`
}

type QueryCreateResp struct {
	Id int64 `json:"id"`
}
