package shared

import (
	"errors"
)

var (
	ErrFramer    = errors.New("error while framer processing")
	ErrDetection = errors.New("error while detection processing")
)
