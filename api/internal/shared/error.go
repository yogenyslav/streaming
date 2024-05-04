package shared

import (
	"errors"
)

var (
	ErrCancelQuery = errors.New("failed to cancel query processing")
)
