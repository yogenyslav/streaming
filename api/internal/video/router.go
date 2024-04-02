package video

import (
	"github.com/gofiber/fiber/v2"
)

type queryHandler interface {
	File(ctx *fiber.Ctx) error
	Link(ctx *fiber.Ctx) error
	CancelById(ctx *fiber.Ctx) error
}

type responseHandler interface {
	FindOneByQueryId(ctx *fiber.Ctx) error
}

func SetupQueryRoutes(app *fiber.App, h queryHandler) {
	g := app.Group("/api/video")

	g.Post("/file", h.File)
	g.Post("/link", h.Link)
	g.Post("/cancel/:id", h.CancelById)
}

func SetupResponseRoutes(app *fiber.App, h responseHandler) {
	g := app.Group("/api/video")

	g.Get("/result/:id", h.FindOneByQueryId)
}
