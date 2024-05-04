package query

import (
	"github.com/gofiber/fiber/v2"
)

type queryHandler interface {
	InsertVideo(c *fiber.Ctx) error
	InsertStream(c *fiber.Ctx) error
	CancelById(c *fiber.Ctx) error
}

func SetupQueryRoutes(app *fiber.App, h queryHandler) {
	g := app.Group("/api/streaming/query")

	g.Post("/video", h.InsertVideo)
	g.Post("/stream", h.InsertStream)
	g.Post("/cancel/:id", h.CancelById)
}
