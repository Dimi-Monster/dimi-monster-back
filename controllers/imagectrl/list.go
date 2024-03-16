package imagectrl

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"pentag.kr/dimimonster/middleware"
	"pentag.kr/dimimonster/models"
)

type ImageResponse struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created-at"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Like        int    `json:"like"`
	LikedByMe   bool   `json:"liked-by-me"`
	Thumbnail   string `json:"thumbnail"`
}

func ListImageCtrl(c *fiber.Ctx) error {
	pageStr := c.Query("page", "0")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		return c.Status(400).SendString("Bad Request")
	}
	userID := middleware.GetUserIDFromMiddleware(c)
	coll := mgm.Coll(&models.Image{})
	// order by created_at desc
	limit := int64(21)
	skip := limit * page
	// offset 0
	images := []models.Image{}

	err = coll.SimpleFind(&images, bson.M{}, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
		Sort:  bson.M{"created_at": -1},
	})
	if err != nil {
		return c.Status(500).SendString("Internal Server Error")
	}
	result := make([]ImageResponse, len(images))
	for i, image := range images {
		var base64String string
		thumbnailfile, err := os.ReadFile("./data/thumbnail/" + image.ID.Hex() + ".jpg")
		if err != nil {
			fmt.Println(err)
			base64String = ""
		} else {
			base64String = "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(thumbnailfile)
		}
		likedByMe := false
		for _, user := range image.LikedBy {
			if user == userID {
				likedByMe = true
				break
			}
		}

		result[i] = ImageResponse{
			ID:          image.ID.Hex(),
			CreatedAt:   image.CreatedAt.Format(time.RFC3339),
			Description: image.Description,
			Location:    image.Location,
			Like:        len(image.LikedBy),
			LikedByMe:   likedByMe,
			Thumbnail:   base64String,
		}
	}
	return c.JSON(result)
}
