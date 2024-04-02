package shared

import (
	"errors"
)

// 400
var (
	ErrQueryIsNotProcessed       = errors.New("query with this id is not being processed")
	ErrUnspecifiedTimeoutForLink = errors.New("queries of type link must have specified timeout in seconds")
)

// 500
var (
	ErrInsertRecord  = errors.New("failed to insert record")
	ErrFindRecord    = errors.New("failed to find record")
	ErrUpdateRecord  = errors.New("failed to update record")
	ErrH264Encode    = errors.New("unable to encode video into H264 format")
	ErrProcessQuery  = errors.New("error while processing query")
	ErrFindProcessed = errors.New("failed to find result of processed query")
)
