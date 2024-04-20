package response

import (
	"net/http"

	"streaming/api/internal/shared"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

var errStatus = map[error]ErrorResponse{
	pgx.ErrNoRows: {
		Msg:    "no rows were found",
		Status: http.StatusNotFound,
	},
	fiber.ErrUnprocessableEntity: {
		Msg:    "data validation error",
		Status: http.StatusUnprocessableEntity,
	},

	// 400
	shared.ErrQueryIsNotProcessed: {
		Status: http.StatusBadRequest,
	},
	shared.ErrUnspecifiedTimeoutForLink: {
		Status: http.StatusBadRequest,
	},

	// 500
	shared.ErrInsertRecord: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrFindRecord: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrUpdateRecord: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrH264Encode: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrProcessQuery: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrFindProcessed: {
		Status: http.StatusInternalServerError,
	},
}
