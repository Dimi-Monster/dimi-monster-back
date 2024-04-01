package imagectrl

import (
	"context"
	"encoding/base64"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"pentag.kr/dimimonster/middleware"
	"pentag.kr/dimimonster/models"
)

func GetWeeklyBestCtrl(c *fiber.Ctx) error {
	needThumbnail := c.Query("thumbnail", "true") == "true"
	userID := middleware.GetUserIDFromMiddleware(c)
	type weeklyAggregate struct {
		HexID       string    `bson:"_id"`
		CreatedAt   time.Time `bson:"created_at"`
		Size        int       `bson:"size"`
		Description string    `bson:"description"`
		Location    string    `bson:"location"`
		LikedBy     []string  `bson:"liked_by"`
	}

	pipeline := bson.A{
		bson.M{
			"$project": bson.M{
				"created_at":  1,
				"description": 1,
				"location":    1,
				"liked_by":    1,
				"size":        bson.M{"$size": "$liked_by"},
			},
		},
		bson.M{
			"$match": bson.M{
				"created_at": bson.M{"$gt": primitive.NewDateTimeFromTime(time.Now().AddDate(0, 0, -7))},
			},
		},
		bson.M{
			"$sort": bson.M{
				"size": -1,
			},
		},
		bson.M{
			"$limit": 3,
		},
	}

	cursor, err := mgm.Coll(&models.Image{}).Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Error(err)
		return c.Status(500).SendString("Internal Server Error")
	}

	var images []weeklyAggregate
	if err = cursor.All(context.Background(), &images); err != nil {
		log.Error(err)
		return c.Status(500).SendString("Internal Server Error")
	}
	var result []ImageResponse
	if needThumbnail {
		for _, image := range images {
			var base64String string
			thumbnailfile, err := os.ReadFile("./data/thumbnail/" + image.HexID + ".jpg")
			if err != nil {
				log.Error(err)
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
			result = append(result, ImageResponse{
				ID:          image.HexID,
				CreatedAt:   image.CreatedAt.Format(time.RFC3339),
				Description: image.Description,
				Location:    image.Location,
				Like:        image.Size,
				LikedByMe:   likedByMe,
				Thumbnail:   base64String,
			})
		}
	} else {
		for _, image := range images {
			likedByMe := false
			for _, user := range image.LikedBy {
				if user == userID {
					likedByMe = true
					break
				}
			}
			result = append(result, ImageResponse{
				ID:          image.HexID,
				CreatedAt:   image.CreatedAt.Format(time.RFC3339),
				Description: image.Description,
				Location:    image.Location,
				Like:        image.Size,
				LikedByMe:   likedByMe,
			})
		}
	}

	return c.JSON(result)
}
