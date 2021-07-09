package views

import (
	"github.com/ewenquim/horkruxes/service"
	"github.com/gofiber/fiber/v2"
)

func SetupLocalRoutes(s service.Service, app *fiber.App) {
	app.Get("/ping", pong)
	// Main view with "Filters"
	app.Get("/", GetMain(s))
	app.Post("/", GetMain(s)) // Select pods by applying filter. Changes the url that Get will parse
	app.Get("/user/:pubKey", GetAuthor(s))
	app.Get("/comments/:uuid", GetComments(s))
	app.Post("/new", NewMessage(s))
	// Keys
	app.Get("/keys", GetKeys)
	app.Post("/keys", PostKeys)
	// FAQ
	app.Get("/faq", GetFaq)
}

// Healthcheck
func pong(c *fiber.Ctx) error {
	return c.SendString("pong")
}
