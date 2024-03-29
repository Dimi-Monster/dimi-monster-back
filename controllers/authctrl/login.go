package authctrl

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
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
		log.Println(err)
		return c.Status(500).SendString("Internal Server Error")
	} else if !profile.EmailVerified || !strings.HasSuffix(profile.Email, "@dimigo.hs.kr") {
		return c.Status(401).SendString("Invalid Email Address")
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
				log.Println("Error while creating user")
				log.Println(err)
				return c.Status(500).SendString("Internal Server Error")
			}
			foundUser = newUser
		} else {
			log.Println("Error while finding user")
			log.Println(err)
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
		log.Println("Error while updating user")
		log.Println(err)
		return c.Status(500).SendString("Internal Server Error")
	}
	accessToken, err := crypt.CreateJWT(crypt.JWTClaims{UserID: foundUser.ID.Hex(), Type: "auth"})
	if err != nil {
		log.Println("Error while creating JWT")
		log.Println(err)
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
