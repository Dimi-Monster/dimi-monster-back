package routers

import (
	"github.com/gofiber/fiber/v2"
)

func InitRouter(app *fiber.App) {
	apiRoter := app.Group("/api")
	initAuthRouter(apiRoter)
	initImageRouter(apiRoter)

	apiRoter.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
}
