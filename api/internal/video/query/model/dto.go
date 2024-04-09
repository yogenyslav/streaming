package model

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type QueryCreateReq struct {
	Source  string                `json:"source" validate:"required"`
	Timeout int64                 `json:"timeout"` // in seconds
	Name    uuid.UUID             `json:"-"`
	File    *multipart.FileHeader `json:"-"`
}

type QueryCreateResp struct {
	Id int64 `json:"id"`
}
