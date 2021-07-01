package api

import (
	"github.com/ewenquim/horkruxes/service"
	"github.com/gofiber/fiber/v2"
)

func SetupApiRoutes(s service.Service, app *fiber.App) {
	app.Get("/api/message", GetMessagesJSON(s))
	app.Get("/api/message/:id", GetMessageJSON(s))
	app.Get("/api/message/author/:pubKey", GetMessagesFromAuthorJSON(s))
	app.Post("/api/message", NewMessage(s))
}
