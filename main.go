package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"pentag.kr/dimimonster/database"
	"pentag.kr/dimimonster/routers"
)

func main() {
	database.InitMDB()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "*",
	}))
	routers.InitRouter(app)
	app.Listen(":3000")
}
