package reportctrl

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"pentag.kr/dimimonster/models"
	"pentag.kr/dimimonster/utils/validator"
)

type ReportProcessRequest struct {
	ReportID string `json:"report-id" validate:"required"`
	Process string `json:"process" validate:"required"`
	Reason string `json:"reason" validate:"required"`
	Secret string `json:"secret" validate:"required"`
}

func ProcessReportCtrl(c *fiber.Ctx) error {
	body := new(ReportProcessRequest)
	if errArr := validator.ParseAndValidate(c, body); errArr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": errArr,
		})
	}
	if !validator.IsHex(body.ReportID) {
		return c.Status(400).SendString("Bad Request")
	}

	foundReport := models.Report{}
	err := mgm.Coll(&foundReport).FindByID(body.ReportID, &foundReport)
	if err != nil {
		return c.Status(404).SendString("Not Found")
	}

	if foundReport.Status != "pending" {
		return c.Status(409).SendString("Already Processed")
	} else if foundReport.Secret != body.Secret {
		return c.Status(403).SendString("Forbidden")
	}

	switch body.Process {
	case "delete":
		foundReport.Status = "delete"
		foundImage := models.Image{}
		err := mgm.Coll(&foundImage).FindByID(foundReport.TargetImageID, &foundImage)
		if err != nil {
			return c.Status(404).SendString("Not Found")
		}
		
		err = mgm.Coll(&foundImage).Delete(&foundImage)
		if err != nil {
			return c.Status(500).SendString("Internal Server Error")
		}		
	case "deleteban":
		foundReport.Status = "deleteban"
		foundImage := models.Image{}
		err := mgm.Coll(&foundImage).FindByID(foundReport.TargetImageID, &foundImage)
		if err != nil {
			return c.Status(404).SendString("Not Found")
		}

		err = mgm.Coll(&foundImage).Delete(&foundImage)
		if err != nil {
			return c.Status(500).SendString("Internal Server Error")
		}
		foundUser := models.User{}
		err = mgm.Coll(&foundUser).FindByID(foundImage.OwnerID, &foundUser)
		if err != nil {
			return c.Status(404).SendString("Not Found")
		}
		foundUser.Banned = true
		foundUser.RefreshTokens = []string{}
		err = mgm.Coll(&foundUser).Update(&foundUser)
		if err != nil {
			return c.Status(500).SendString("Internal Server Error")
		}

	case "withdraw":
		foundReport.Status = "withdraw"
	case "reporterban":
		foundReport.Status = "reporterban"
		foundUser := models.User{}
		err := mgm.Coll(&foundUser).FindByID(foundReport.ReporterID, &foundUser)
		if err != nil {
			return c.Status(404).SendString("Not Found")
		}
		foundUser.Banned = true
		foundUser.RefreshTokens = []string{}
		err = mgm.Coll(&foundUser).Update(&foundUser)
		if err != nil {
			return c.Status(500).SendString("Internal Server Error")
		}
	default:
		return c.Status(400).SendString("Bad Request")
	}
	foundReport.ProcessReason = body.Reason
	err = mgm.Coll(&foundReport).Update(&foundReport)
	if err != nil {
		return c.Status(500).SendString("Internal Server Error")
	}
	return c.SendString("Done")
}
