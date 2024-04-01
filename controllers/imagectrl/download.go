package imagectrl

import (
	"github.com/gofiber/fiber/v2"
	"pentag.kr/dimimonster/utils/validator"
)

func DownloadImageCtrl(c *fiber.Ctx) error {
	imageID := c.Params("id")
	if !validator.IsHex(imageID) {
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}
	c.Set("Content-Type", "image/jpeg")
	return c.SendFile("./data/original/" + imageID + ".jpg")
}
