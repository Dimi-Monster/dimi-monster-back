package crypt

import (
	"log"

	"github.com/go-resty/resty/v2"
	"pentag.kr/dimimonster/config"
)

type RecaptchaReqBody struct {
	Event RecaptchaEvent `json:"event"`
}

type RecaptchaEvent struct {
	Token          string `json:"token"`
	SiteKey        string `json:"siteKey"`
	ExpectedAction string `json:"expectedAction"`
}

type RecaptchaResBody struct {
	RiskAnalysis RecaptchaRisk `json:"riskAnalysis"`
}

type RecaptchaRisk struct {
	Score float32 `json:"score"`
}

func RecaptchaCheck(token string, action string) bool {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(RecaptchaReqBody{
			Event: RecaptchaEvent{
				Token:          token,
				SiteKey:        config.RecaptchaSecret,
				ExpectedAction: action,
			},
		},
		).
		SetResult(&RecaptchaResBody{}).
		Post("https://recaptchaenterprise.googleapis.com/v1/projects/dimi-monster/assessments?key=" + config.RecaptchaAPIKey)
	if err != nil {
		log.Println(err)
		return false
	}
	result := resp.Result().(*RecaptchaResBody)
	return result.RiskAnalysis.Score > 0.5
}
