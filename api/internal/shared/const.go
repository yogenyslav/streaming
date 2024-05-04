package shared

const (
	FormDataVideoKey = "video"

	VideoBucket = "video"
)

type QueryType int8

const (
	QueryTypeUndefined QueryType = iota
	QueryTypeVideo
	QueryTypeStream
)

type ResponseStatus int8

const (
	ResponseStatusUndefined ResponseStatus = iota
	ResponseStatusInit
	ResponseStatusFramerProcessing
	ResponseStatusFramerError
	ResponseStatusDetectionProcessing
	ResponseStatusDetectionError
	ResponseStatusSuccess
	ResponseStatusError
	ResponseStatusCanceled
	ResponseStatusResponserProcessing
)
