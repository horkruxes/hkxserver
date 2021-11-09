package views

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/horkruxes/hkxserver/model"
	"github.com/horkruxes/hkxserver/service"
)

func FromFormToPayload(c *fiber.Ctx) model.Message {
	message := model.Message{}

	message.SignatureBase64 = strings.TrimSpace(c.FormValue("signature"))
	message.AuthorBase64 = strings.TrimSpace(c.FormValue("public-key"))
	message.Content = strings.TrimSpace(c.FormValue("message"))
	message.DisplayedName = strings.TrimSpace(c.FormValue("name"))
	message.MessageID = strings.TrimSpace(c.FormValue("answer-to"))
	message.Pod = strings.TrimSpace(c.FormValue("pod-to-post-to"))

	return message
}

func parsePodsListToURL(s service.Service, text string) string {
	list := s.Regexes.URLs.FindAllString(text, -1)
	text = strings.Join(list, "+")
	return strings.ToLower(text)
}

func parseFormsToService(c *fiber.Ctx, s service.Service) service.ClientConfig {
	public := !(c.FormValue("hide-friends") == "on")
	private := c.FormValue("pods") != ""
	list := c.FormValue("pods")

	s.ClientConfig.PublicPods = public
	s.ClientConfig.SpecificPods = private
	s.ClientConfig.SpecificPodsListString = parsePodsListToURL(s, list)
	s.ClientConfig.SpecificPodsList = strings.Split(s.ClientConfig.SpecificPodsListString, "+")
	// fmt.Println("s list", s.ClientConfig.SpecificPodsList)

	s.UpdateClientPodsList()
	return s.ClientConfig
}
