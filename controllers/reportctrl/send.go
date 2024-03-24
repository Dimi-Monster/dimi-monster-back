package reportctrl

import "github.com/gofiber/fiber/v2"

type ReportRequest struct {
	TargetID string `json:"target-id" validate:"required"`
	Category int    `json:"category" validate:"required"`
	Reason   string `json:"reason" validate:"required"`
	Token    string `json:"token" validate:"required"`
}

func SendReportCtrl(c *fiber.Ctx) error {

	return c.SendString("Send Report")
}
