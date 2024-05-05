package response

import (
	"github.com/gofiber/fiber/v2"
)

type responseHandler interface {
	GetResult(c *fiber.Ctx) error
}

func SetupResponseRoutes(app *fiber.App, h responseHandler) {
	g := app.Group("/api/streaming/response")

	g.Get("/:id", h.GetResult)
}
