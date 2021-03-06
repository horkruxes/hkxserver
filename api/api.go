package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/horkruxes/hkxserver/service"
)

func SetupApiRoutes(s service.Service, app *fiber.App) {
	if s.ServerConfig.Enabled {
		app.Get("/api/all", GetAllJSON(s))
		app.Get("/api/message", GetMessagesJSON(s))
		app.Get("/api/message/:id", GetMessageJSON(s))
		app.Get("/api/comments/:id", GetCommentsJSON(s))
		app.Get("/api/user/:pubKey", GetMessagesFromAuthorJSON(s))
		app.Get("/api/server-info", func(c *fiber.Ctx) error {
			return c.Status(200).JSON(s.GeneralConfig)
		})
		if !s.ServerConfig.LockedByDefault {
			app.Post("/api/message", NewMessage(s))
		}
	}
}
