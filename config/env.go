package config

import (
	"fmt"
	"os"
)

var JWTSecret = []byte(LoadEnv("JWT_SECRET"))
var RecaptchaSecret = LoadEnv("RECAPTCHA_SECRET")
var RecaptchaAPIKey = LoadEnv("RECAPTCHA_API_KEY")
var DB_URI = LoadEnv("DB_URI")
var GAUTH_SECRET = LoadEnv("GAUTH_SECRET")

func LoadEnv(key string) (result string) {
	result = os.Getenv(key)
	if result == "" {
		panic(fmt.Errorf("env %s is empty", key))
	}
	fmt.Println(key, ":", result)
	return
}
