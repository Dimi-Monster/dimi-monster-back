package authctrl

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"pentag.kr/dimimonster/config"
	"pentag.kr/dimimonster/models"
	"pentag.kr/dimimonster/utils/crypt"
	"pentag.kr/dimimonster/utils/gauth"
	"pentag.kr/dimimonster/utils/random"
)

func LoginCtrl(c *fiber.Ctx) error {
	code := c.Query("code", "")
	if code == "" {
		return c.Status(400).SendString("code is required")
	}
	profile, err := gauth.GetProfile(code)
	if err != nil {
		log.Error(err)
		return c.Status(500).SendString("Internal Server Error")
	} else if !profile.EmailVerified {
		return c.Status(401).SendString("Invalid Email Address")
	} else if !strings.HasSuffix(profile.Email, "@dimigo.hs.kr") {
		foundEmailWhite := &models.EmailWhitelist{}
		err = mgm.Coll(foundEmailWhite).First(bson.M{"email": profile.Email}, foundEmailWhite)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(401).SendString("Invalid Email Address")
			} else {
				log.Error("Error while finding email whitelist")
				log.Error(err)
				return c.Status(500).SendString("Internal Server Error")
			}
		}
	}
	foundUser := &models.User{}
	err = mgm.Coll(foundUser).First(bson.M{"gid": profile.SUB}, foundUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			newUser := models.NewUser(
				profile.SUB,
				profile.Name,
				profile.Email,
			)
			err = mgm.Coll(newUser).Create(newUser)
			if err != nil {
				log.Error("Error while creating user")
				log.Error(err)
				return c.Status(500).SendString("Internal Server Error")
			}
			foundUser = newUser
		} else {
			log.Error("Error while finding user")
			log.Error(err)
			return c.Status(500).SendString("Internal Server Error")
		}
	}
	if foundUser.Banned {
		return c.Status(403).SendString("Banned User")
	}
	if len(foundUser.RefreshTokens) > 6 {
		foundUser.RefreshTokens = foundUser.RefreshTokens[len(foundUser.RefreshTokens)-5:]
	}
	foundUser.RefreshTokens = append(foundUser.RefreshTokens, random.RandString(32))
	err = mgm.Coll(foundUser).Update(foundUser)
	if err != nil {
		log.Error("Error while updating user")
		log.Error(err)
		return c.Status(500).SendString("Internal Server Error")
	}
	accessToken, err := crypt.CreateJWT(crypt.JWTClaims{UserID: foundUser.ID.Hex(), Type: "auth"})
	if err != nil {
		log.Error("Error while creating JWT")
		log.Error(err)
		return c.Status(500).SendString("Internal Server Error")
	}

	return c.JSON(fiber.Map{
		"name":          foundUser.Name,
		"email":         foundUser.Email,
		"picture":       profile.Picture,
		"locale":        profile.Locale,
		"access-token":  accessToken,
		"at-expire":     config.AccessTokenExpire,
		"refresh-token": foundUser.RefreshTokens[len(foundUser.RefreshTokens)-1],
	})
}
