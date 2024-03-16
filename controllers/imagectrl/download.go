package imagectrl

import (
	"github.com/gofiber/fiber/v2"
	"pentag.kr/dimimonster/utils/crypt"
)

func DownloadImageCtrl(c *fiber.Ctx) error {
	imageID := c.Params("id")
	if !isHex(imageID) {
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}
	imageToken := c.Query("image-token", "")
	claims, err := crypt.ValidateJWT(imageToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	} else if claims.Type != "image" {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
	c.Set("Content-Type", "image/jpeg")
	return c.SendFile("./data/original/" + imageID + ".jpg")
}

func isHex(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}
