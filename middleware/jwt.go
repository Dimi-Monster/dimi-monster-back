package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"pentag.kr/dimimonster/utils/crypt"
)

func AuthMiddleware(c *fiber.Ctx) error {
	// get Bearer authHeader
	authHeader := strings.Split(c.Get("Authorization"), " ")
	if len(authHeader) != 2 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	token := authHeader[1]

	claims, err := crypt.ValidateJWT(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	if claims.Type != "auth" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	c.Locals("user-id", claims.UserID) // type: uuid.UUID
	return c.Next()
}

func GetUserIDFromMiddleware(c *fiber.Ctx) string {
	return c.Locals("user-id").(string)
}
