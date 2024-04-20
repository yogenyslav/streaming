package model

type RawFrame struct {
	Id      int
	QueryId int64
	Data    []byte
	IsLast  bool
	Error   error
}

type ProcessedFrame struct {
	Id      int
	QueryId int64
	Source  string
	IsLast  bool
	Error   error
}
