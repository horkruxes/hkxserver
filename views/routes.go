package views

import (
	"github.com/gofiber/fiber/v2"
	"github.com/horkruxes/hkxserver/service"
)

func SetupLocalRoutes(s service.Service, app *fiber.App) {
	app.Get("/ping", pong)
	if s.ClientConfig.Enabled {
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
		app.Get("/404", page404)
	} else {
		app.Get("/", GetMainNoFront)
	}
}

// Healthcheck
func pong(c *fiber.Ctx) error {
	return c.SendString("pong")
}

func page404(c *fiber.Ctx) error {
	return c.SendString("404: wrong url")
}
