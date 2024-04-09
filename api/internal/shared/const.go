package shared

type QueryType int32

const (
	TypeFile QueryType = iota
	TypeLink
)

type ResponseStatus int32

const (
	StatusSuccess ResponseStatus = iota
	StatusProcessing
	StatusError
	StatusCanceled
)

const (
	UniqueViolationCode = "23505"
	PreprocessedFileExt = ".mp4"
	VideoBucket         = "video"
)
