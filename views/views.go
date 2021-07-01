package views

import (
	"github.com/ewenquim/horkruxes/service"
	"github.com/fatih/structs"
	"github.com/gofiber/fiber/v2"
)

func SetupLocalRoutes(s service.Service, app *fiber.App) {
	app.Get("/ping", pong)
	app.Get("/", GetMain(s))
	app.Get("/keys", GetKeys)
	// app.Get("/api/message/author/:pubKey", GetMessagesFromAuthorJSON(s.DB))
	app.Get("/pubkey/:pubKey", GetAuthor(s))
	app.Post("/keys", PostKeys)
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
		id :=  SafeURLToBase64(c.Params("pubKey"))
		localData := GetAuthorMessagesAndMainPageInfo(s, id)
		return c.Render("main/root", structs.Map(localData))
	}
}
