package server

import (
	"net/http"

	"streaming/api/internal/shared"

	"github.com/yogenyslav/pkg/response"
)

var errStatus = map[error]response.ErrorResponse{
	// 400
	shared.ErrCancelQuery: {
		Status: http.StatusBadRequest,
	},
}
