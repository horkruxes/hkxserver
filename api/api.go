package api

import (
	"github.com/ewenquim/horkruxes-client/model"
	"github.com/gofiber/fiber/v2"
)

func SetupMessagesRoutes(s model.Service, app *fiber.App) {
	app.Get("/api/message", s.GetMessagesJSON)
	app.Get("/api/message/:id", s.GetMessage)
	app.Post("/api/message", s.NewMessage)
	app.Delete("/api/message/:id", s.DeleteMessage)
}
