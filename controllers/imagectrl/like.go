package imagectrl

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"pentag.kr/dimimonster/middleware"
	"pentag.kr/dimimonster/models"
	"pentag.kr/dimimonster/utils/validator"
)

func GetLikeCountCtrl(c *fiber.Ctx) error {
	imageID := c.Params("id")
	if !validator.IsHex(imageID) {
		return c.Status(400).SendString("Bad Request")
	}
	image := models.Image{}
	err := mgm.Coll(&models.Image{}).FindByID(imageID, &image)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).SendString("Not Found")
		}
		return c.Status(500).SendString("Internal Server Error")
	}
	return c.JSON(fiber.Map{
		"like": len(image.LikedBy),
	})
}

func AddLikeCtrl(c *fiber.Ctx) error {
	imageID := c.Params("id")
	if !validator.IsHex(imageID) {
		return c.Status(400).SendString("Bad Request")
	}
	userID := middleware.GetUserIDFromMiddleware(c)
	image := models.Image{}
	err := mgm.Coll(&models.Image{}).FindByID(imageID, &image)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).SendString("Not Found")
		}
		return c.Status(500).SendString("Internal Server Error")
	}
	if image.HasLike(userID) {
		return c.Status(400).SendString("Bad Request")
	}
	image.AddLike(userID)
	err = mgm.Coll(&models.Image{}).Update(&image)
	if err != nil {
		return c.Status(500).SendString("Internal Server Error")
	}

	return c.SendString("Done")
}

func DeleteLikeCtrl(c *fiber.Ctx) error {
	imageID := c.Params("id")
	if !validator.IsHex(imageID) {
		return c.Status(400).SendString("Bad Request")
	}
	userID := middleware.GetUserIDFromMiddleware(c)
	image := models.Image{}
	err := mgm.Coll(&models.Image{}).FindByID(imageID, &image)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).SendString("Not Found")
		}
		return c.Status(500).SendString("Internal Server Error")
	}
	if !image.HasLike(userID) {
		return c.Status(400).SendString("Bad Request")
	}
	image.DeleteLike(userID)
	err = mgm.Coll(&models.Image{}).Update(&image)
	if err != nil {
		return c.Status(500).SendString("Internal Server Error")
	}
	return c.SendString("Done")
}
