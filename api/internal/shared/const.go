package shared

type QueryType int32

const (
	QueryTypeUnspecified QueryType = iota
	QueryTypeVideo       QueryType = iota
	QueryTypeStream
)

type ResponseStatus int32

const (
	ResponseStatusUnspecified ResponseStatus = iota
	ResponseStatusInitProcessing
	ResponseStatusFramerProcessing
	ResponseStatusFramerSuccess
	ResponseStatusFramerError
	ResponseStatusDetectionProcessing
	ResponseStatusDetectionSuccess
	ResponseStatusDetectionError
	ResponseStatusSuccess
	ResponseStatusCanceled
	ResponseStatusError
)

const (
	UniqueViolationCode = "23505"
	VideoBucket         = "video"
)
