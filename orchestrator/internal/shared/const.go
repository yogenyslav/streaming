package shared

const (
	FramerTopic          = "framer"
	FramerResTopic       = "framer-res"
	DetectionCancelTopic = "detection-cancel"
	DetectionResTopic    = "detection-res"
	ResponserResTopic    = "responser-res"

	SuccessMessage = "success"
	ErrorMessage   = "error"
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
