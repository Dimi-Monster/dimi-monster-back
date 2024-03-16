package imagectrl

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"pentag.kr/dimimonster/models"
)

func GetRecentImageCtrl(c *fiber.Ctx) error {

	images := []models.Image{}
	err := mgm.Coll(&models.Image{}).SimpleFind(&images, &bson.M{}, options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(1))
	if err != nil {
		log.Println(err)
		return c.Status(500).SendString("Internal Server Error")
	}

	return c.JSON(
		fiber.Map{
			"id": images[0].ID.Hex(),
		},
	)
}
