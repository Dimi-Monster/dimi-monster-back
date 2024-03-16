package crypt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"pentag.kr/dimimonster/config"
)

type JWTClaims struct {
	UserID string
	Type   string
}

func CreateJWT(claimData JWTClaims) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = claimData.UserID
	claims["type"] = claimData.Type
	claims["exp"] = time.Now().Add(config.AccessTokenExpire * time.Second).Unix()
	return token.SignedString(config.JWTSecret)
}

func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.JWTSecret, nil
	})
	if err != nil {
		return jwt.MapClaims{}, err
	} else if token.Claims == nil {
		return jwt.MapClaims{}, err
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
}

func ValidateJWT(tokenString string) (*JWTClaims, error) {
	claims, err := ParseJWT(tokenString)
	if err != nil {
		return nil, err
	}
	return &JWTClaims{
		UserID: claims["user_id"].(string),
		Type:   claims["type"].(string),
	}, nil
}
