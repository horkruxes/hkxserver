package api

import (
	"github.com/ewenquim/horkruxes-client/service"
	"github.com/gofiber/fiber/v2"
)

func SetupApiRoutes(s service.Service, app *fiber.App) {
	app.Get("/api/message", GetMessagesJSON(s.DB))
	app.Get("/api/message/:id", GetMessageJSON(s.DB))
	app.Post("/api/message", NewMessage(s.DB))
}
