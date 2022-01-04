package views

import (
	"html"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"github.com/horkruxes/hkxserver/service"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

func parseFormsToService(c *fiber.Ctx, s service.Service) service.ClientConfig {
	public := !(c.FormValue("hide-friends") == "on")

	s.ClientConfig.PublicPods = public

	return s.ClientConfig
}

func MarkDowner(content string) template.HTML {
	// Unescapes from db bacause blackfriday will escape them
	content = html.UnescapeString(content)

	markdownBytes := blackfriday.Run(
		[]byte(content),
		blackfriday.WithExtensions(blackfriday.HardLineBreak|blackfriday.NoEmptyLineBeforeBlock),
	)

	// Clean Client-side XSS
	markdownBytes = bluemonday.StrictPolicy().SanitizeBytes(markdownBytes)

	//#nosec
	return template.HTML(markdownBytes)
}
