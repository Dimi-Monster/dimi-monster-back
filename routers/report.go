package routers

import (
	"github.com/gofiber/fiber/v2"
	"pentag.kr/dimimonster/controllers/reportctrl"
	"pentag.kr/dimimonster/middleware"
)

func initReportRouter(router fiber.Router) {
	reportRouter := router.Group("/report")
	reportRouter.Use(middleware.AuthMiddleware)

	reportRouter.Post("/",
		func(c *fiber.Ctx) error {
			return reportctrl.SendReportCtrl(c)
		},
	)
}
