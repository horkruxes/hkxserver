package views

import (
	"github.com/ewenquim/horkruxes/service"
	"github.com/fatih/structs"
	"github.com/gofiber/fiber/v2"
)

func SetupLocalRoutes(s service.Service, app *fiber.App) {
	app.Get("/ping", pong)
	// Main view with "Filters"
	app.Get("/", GetMain(s))
	app.Get("/user/:pubKey", GetAuthor(s))
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

func GetMain(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		localData := GetMessagesAndMainPageInfo(s)
		return c.Render("main/root", structs.Map(localData))
	}
}

func GetAuthor(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := SafeURLToBase64(c.Params("pubKey"))
		localData := GetAuthorMessagesAndMainPageInfo(s, id)
		return c.Render("main/root", structs.Map(localData))
	}
}
