package views

import (
	"html/template"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/horkruxes/hkxserver/model"
	"github.com/horkruxes/hkxserver/service"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

func FormToBasicMessage(c *fiber.Ctx) model.Message {
	message := model.Message{}

	message.SignatureBase64 = strings.TrimSpace(c.FormValue("signature"))
	message.AuthorBase64 = strings.TrimSpace(c.FormValue("public-key"))
	message.Content = strings.TrimSpace(c.FormValue("message"))
	message.DisplayedName = strings.TrimSpace(c.FormValue("name"))
	message.MessageID = strings.TrimSpace(c.FormValue("answer-to"))
	message.Pod = strings.TrimSpace(c.FormValue("pod-to-post-to"))

	return message
}

func parseFormsToService(c *fiber.Ctx, s service.Service) service.ClientConfig {
	public := !(c.FormValue("hide-friends") == "on")

	s.ClientConfig.PublicPods = public

	return s.ClientConfig
}

func MarkDowner(policy *bluemonday.Policy) func(string) template.HTML {
	return func(content string) template.HTML {
		markdownBytes := blackfriday.Run([]byte(content), blackfriday.WithExtensions(blackfriday.HardLineBreak|blackfriday.NoEmptyLineBeforeBlock))
		safeBytes := policy.SanitizeBytes(markdownBytes)
		return template.HTML(safeBytes)
	}
}
