package authctrl

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"pentag.kr/dimimonster/models"
	"pentag.kr/dimimonster/utils/validator"
)

func LogoutCtrl(c *fiber.Ctx) error {
	body := new(RefreshRequest)
	if errArr := validator.ParseAndValidate(c, body); errArr != nil {
		return c.Status(400).JSON(fiber.Map{
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
		return c.Status(404).SendString("Refresh Token Not Found")
	}
	foundUser.RemoveRefreshToken(body.RefreshToken)
	err = mgm.Coll(foundUser).Update(foundUser)
	if err != nil {
		return c.Status(500).SendString("Internal Server Error")
	}
	return c.SendString("Logout")
}
