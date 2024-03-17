package imagectrl

import (
	"github.com/gofiber/fiber/v2"
)

func DownloadImageCtrl(c *fiber.Ctx) error {
	imageID := c.Params("id")
	if !isHex(imageID) {
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
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
