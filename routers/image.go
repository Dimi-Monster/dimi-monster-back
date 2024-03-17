package routers

import (
	"github.com/gofiber/fiber/v2"
	"pentag.kr/dimimonster/controllers/imagectrl"
	"pentag.kr/dimimonster/middleware"
)

func initImageRouter(router fiber.Router) {
	imageRouter := router.Group("/image")
	imageRouter.Use(middleware.AuthMiddleware)

	imageRouter.Get("/recent",
		func(c *fiber.Ctx) error {
			return imagectrl.GetRecentImageCtrl(c)
		},
	)

	imageRouter.Get("/list",
		func(c *fiber.Ctx) error {
			return imagectrl.ListImageCtrl(c)
		},
	)
	imageRouter.Post("/upload",
		func(c *fiber.Ctx) error {
			return imagectrl.UploadImage(c)
		},
	)

	imageRouter.Get("/like/:id",
		func(c *fiber.Ctx) error {
			return imagectrl.GetLikeCountCtrl(c)
		},
	)
	imageRouter.Post("/like/:id",
		func(c *fiber.Ctx) error {
			return imagectrl.AddLikeCtrl(c)
		},
	)
	imageRouter.Delete("/like/:id",
		func(c *fiber.Ctx) error {
			return imagectrl.DeleteLikeCtrl(c)
		},
	)

	imageRouter.Get("/weekly",
		func(c *fiber.Ctx) error {
			return imagectrl.GetWeeklyBestCtrl(c)
		},
	)

	imageRouter.Get("/:id",
		func(c *fiber.Ctx) error {
			return imagectrl.DownloadImageCtrl(c)
		},
	)
}
