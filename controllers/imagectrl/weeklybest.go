package imagectrl

import (
	"encoding/base64"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"pentag.kr/dimimonster/middleware"
	"pentag.kr/dimimonster/models"
)

func GetWeeklyBestCtrl(c *fiber.Ctx) error {
	images := []models.Image{}
	// get recent 1 week images
	// which has most longest liked_by length
	limit := int64(1)
	err := mgm.Coll(&models.Image{}).SimpleFind(&images, &bson.M{
		"created_at": bson.M{"$gte": time.Now().AddDate(0, 0, -7)},
	}, &options.FindOptions{
		Sort:  bson.M{"liked_by": -1},
		Limit: &limit,
	})
	if err != nil {
		log.Panicln(err)
		return c.Status(500).SendString("Internal Server Error")
	}
	var base64String string
	thumbnailfile, err := os.ReadFile("./data/thumbnail/" + images[0].ID.Hex() + ".jpg")
	if err != nil {
		log.Panicln(err)
		base64String = ""
	} else {
		base64String = "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(thumbnailfile)
	}

	return c.JSON(
		ImageResponse{
			ID:          images[0].ID.Hex(),
			CreatedAt:   images[0].CreatedAt.Format(time.RFC3339),
			Description: images[0].Description,
			Location:    images[0].Location,
			Like:        len(images[0].LikedBy),
			LikedByMe:   images[0].HasLike(middleware.GetUserIDFromMiddleware(c)),
			Thumbnail:   base64String,
		},
	)
}
