package views

import (
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

func MarkDowner(policy *bluemonday.Policy) func(string) template.HTML {
	return func(content string) template.HTML {
		markdownBytes := blackfriday.Run(
			[]byte(content),
			blackfriday.WithExtensions(blackfriday.HardLineBreak|blackfriday.NoEmptyLineBeforeBlock),
		)
		//#nosec
		return template.HTML(markdownBytes)
	}
}
