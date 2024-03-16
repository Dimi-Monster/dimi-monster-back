package routers

import (
	"github.com/gofiber/fiber/v2"
	"pentag.kr/dimimonster/controllers/authctrl"
)

func initAuthRouter(router fiber.Router) {
	authRouter := router.Group("/auth")
	authRouter.Get("/login",
		func(c *fiber.Ctx) error {
			return authctrl.LoginCtrl(c)
		},
	)
	authRouter.Post("/refresh",
		func(c *fiber.Ctx) error {
			return authctrl.RefreshCtrl(c)
		},
	)
	authRouter.Delete("/logout",
		func(c *fiber.Ctx) error {
			return authctrl.LogoutCtrl(c)
		},
	)
}
