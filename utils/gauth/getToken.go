package gauth

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"pentag.kr/dimimonster/config"
)

type GoogleTokenPayload struct {
	AccessToken string `json:"access_token"`
}

func GetProfile(code string) (*GooglePayload, error) {
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"code":          code,
			"client_id":     "752669101446-ssmoaio24ohfv2vhg59gphbbdqtpe7kb.apps.googleusercontent.com",
			"client_secret": config.GAUTH_SECRET,
			"redirect_uri":  "https://dimi.monster/redirect/gauth",
			"grant_type":    "authorization_code",
		}).
		Post("https://oauth2.googleapis.com/token")
	if err != nil {
		return nil, err
	}
	payload := GoogleTokenPayload{}
	err = json.Unmarshal(resp.Body(), &payload)
	if err != nil {
		return nil, err
	}
	profile, err := convertToken(payload.AccessToken)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
