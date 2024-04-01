package discord

import (
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/gtuk/discordwebhook"
	"pentag.kr/dimimonster/config"
	"pentag.kr/dimimonster/models"
)

var username = "신고봇"
var color = "16744272"

func SendReportWebhook(report *models.Report) {
	links := "[이미지 보기](https://dimi.monster/?id=&&id&&) [이미지 삭제 처리](https://dimi.monster/admin/delete/&&id&&?key=&&key&&)"
	links = strings.ReplaceAll(links, "&&id&&", report.TargetImageID)
	links = strings.ReplaceAll(links, "&&key&&", report.Secret)
	embed := discordwebhook.Embed{
		Title: &report.Category,
		Description: &report.Reason,
		Color: &color,
	}

	message := discordwebhook.Message{
		Username: &username,
		Embeds: &[]discordwebhook.Embed{embed},
		Content: &links,
	}

	err := discordwebhook.SendMessage(config.DiscordWebhookURL, message)
	if err != nil {
		log.Error(err)
	}
}