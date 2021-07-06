package views

import (
	"strings"

	"github.com/ewenquim/horkruxes/api"
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
