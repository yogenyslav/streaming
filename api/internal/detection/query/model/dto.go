package model

import (
	"io"

	"streaming/api/internal/shared"
)

type StreamCreateReq struct {
	Source string `json:"source"`
}

type QueryCreateReq struct {
	Source string
	Video  io.Reader
	Type   shared.QueryType
	Name   string
	Size   int64
}

type QueryCreateResp struct {
	Id int64 `json:"id"`
}
