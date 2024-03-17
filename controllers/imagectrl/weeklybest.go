package imagectrl

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"pentag.kr/dimimonster/middleware"
	"pentag.kr/dimimonster/models"
)

func GetWeeklyBestCtrl(c *fiber.Ctx) error {
	userID := middleware.GetUserIDFromMiddleware(c)
	type result struct {
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
			"$limit": 1,
		},
	}

	cursor, err := mgm.Coll(&models.Image{}).Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Println(err)
		return c.Status(500).SendString("Internal Server Error")
	}

	var image result
	if !cursor.Next(context.Background()) {
		return c.Status(404).SendString("Not Found")
	}
	err = cursor.Decode(&image)
	if err != nil {
		log.Panicln(err)
		return c.Status(500).SendString("Internal Server Error")
	}

	var base64String string
	thumbnailfile, err := os.ReadFile("./data/thumbnail/" + image.HexID + ".jpg")
	if err != nil {
		log.Panicln(err)
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

	return c.JSON(
		ImageResponse{
			ID:          image.HexID,
			CreatedAt:   image.CreatedAt.Format(time.RFC3339),
			Description: image.Description,
			Location:    image.Location,
			Like:        image.Size,
			LikedByMe:   likedByMe,
			Thumbnail:   base64String,
		},
	)
}
