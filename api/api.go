package api

import (
	"github.com/ewenquim/horkruxes/service"
	"github.com/gofiber/fiber/v2"
)

func SetupApiRoutes(s service.Service, app *fiber.App) {
	app.Get("/api/message", GetMessagesJSON(s))
	app.Get("/api/message/:id", GetMessageJSON(s))
	app.Get("/api/comments/:id", GetCommentsJSON(s))
	app.Get("/api/user/:pubKey", GetMessagesFromAuthorJSON(s))
	app.Get("/api/server-info", GetServerInfo(s))
	app.Post("/api/message", NewMessage(s))
}
