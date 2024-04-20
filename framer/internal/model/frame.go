package model

type RawFrame struct {
	Id      int
	QueryId int64
	Data    []byte
	IsLast  bool
	Error   error
}
