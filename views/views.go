package views

import (
	"github.com/ewenquim/horkruxes-client/service"
	"github.com/fatih/structs"
	"github.com/gofiber/fiber/v2"
)

func SetupLocalRoutes(s service.Service, app *fiber.App) {
	app.Get("/ping", pong)
	app.Get("/", GetMain(s))
	app.Get("/keys", GetKeys)
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
