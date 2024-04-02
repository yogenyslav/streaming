package model

type QueryCreateReq struct {
	Source  string `json:"source" validate:"required"`
	FileExt string `json:"-"`
	Timeout int64  `json:"timeout"` // in seconds
}

type QueryCreateResp struct {
	Id int64 `json:"id"`
}
