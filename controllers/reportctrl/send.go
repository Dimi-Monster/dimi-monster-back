package reportctrl

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/kamva/mgm/v3"
	"pentag.kr/dimimonster/middleware"
	"pentag.kr/dimimonster/models"
	"pentag.kr/dimimonster/utils/validator"
)

type ReportRequest struct {
	TargetID string `json:"target-id" validate:"required"`
	Reason   string `json:"reason" validate:"required, max=300"`
	Token    string `json:"token" validate:"required"`
}

func SendReportCtrl(c *fiber.Ctx) error {
	body := new(ReportRequest)
	if errArr := validator.ParseAndValidate(c, body); errArr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": errArr,
		})
	}
	userID := middleware.GetUserIDFromMiddleware(c)

	if !validator.IsHex(body.TargetID) {
		return c.Status(400).SendString("Bad Request")
	}
	// if !crypt.RecaptchaCheck(body.Token, "image_report") {
	// 	return c.Status(425).SendString("Recaptcha Failed")
	// }

	newReport, err := models.NewReport(
		body.TargetID,
		userID,
		body.Reason,
	)
	if err != nil {
		log.Error(err)	
		return c.Status(500).SendString("Internal Server Error")
	}

	err = mgm.Coll(newReport).Create(newReport)
	if err != nil {
		log.Error(err)
		return c.Status(500).SendString("Internal Server Error")
	}

	return c.SendString("Send Report")
}
