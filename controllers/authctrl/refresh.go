package authctrl

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"pentag.kr/dimimonster/config"
	"pentag.kr/dimimonster/models"
	"pentag.kr/dimimonster/utils/crypt"
	"pentag.kr/dimimonster/utils/validator"
)

type RefreshRequest struct {
	Email        string `json:"email" validate:"required,email"`
	RefreshToken string `json:"refresh-token" validate:"required"`
}

func RefreshCtrl(c *fiber.Ctx) error {
	body := new(RefreshRequest)
	if errArr := validator.ParseAndValidate(c, body); errArr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errArr,
		})
	}
	foundUser := &models.User{}
	err := mgm.Coll(foundUser).First(bson.M{"email": body.Email}, foundUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(401).SendString("Unauthorized")
		}
		return c.Status(500).SendString("Internal Server Error")
	}
	if !foundUser.HasRefreshToken(body.RefreshToken) {
		return c.Status(401).SendString("Unauthorized")
	}
	if foundUser.Banned {
		return c.Status(403).SendString("Banned User")
	}
	accessToken, err := crypt.CreateJWT(crypt.JWTClaims{UserID: foundUser.ID.Hex(), Type: "auth"})
	if err != nil {
		return c.Status(500).SendString("Internal Server Error")
	}
	return c.JSON(fiber.Map{
		"access-token": accessToken,
		"at-expire":    config.AccessTokenExpire,
	})
}
