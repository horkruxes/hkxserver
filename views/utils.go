package views

import (
	"strings"

	"github.com/ewenquim/horkruxes/api"
	"github.com/ewenquim/horkruxes/service"
	"github.com/gofiber/fiber/v2"
)

func FromFormToPayload(c *fiber.Ctx) api.NewMessagePayload {
	message := api.NewMessagePayload{}

	message.Signature = strings.TrimSpace(c.FormValue("signature"))
	message.PublicKey = strings.TrimSpace(c.FormValue("public-key"))
	message.Content = strings.TrimSpace(c.FormValue("message"))
	message.Name = strings.TrimSpace(c.FormValue("name"))
	message.MessageID = strings.TrimSpace(c.FormValue("answer-to"))
	message.Pod = strings.TrimSpace(c.FormValue("pod-to-post-to"))

	return message
}

func parsePodsListToURL(s service.Service, text string) string {
	list := s.Regexes.URLs.FindAllString(text, -1)
	text = strings.Join(list, ";")
	return strings.ToLower(text)
}

func parseFormsToService(c *fiber.Ctx, s service.Service) service.ClientConfig {
	public := c.FormValue("friend-sources", "on") == "on"
	private := c.FormValue("pods") != ""
	list := c.FormValue("pods")

	// fmt.Println("pub", public)
	// fmt.Println("private", private)
	// fmt.Println("pods", list)

	s.ClientConfig.PublicPods = public
	s.ClientConfig.SpecificPods = private
	s.ClientConfig.SpecificPodsListString = parsePodsListToURL(s, list)
	s.ClientConfig.SpecificPodsList = strings.Split(s.ClientConfig.SpecificPodsListString, ";")
	// fmt.Println("s list", s.ClientConfig.SpecificPodsList)

	s.UpdateClientPodsList()
	return s.ClientConfig
}
