package views

import (
	"github.com/ewenquim/horkruxes-client/service"
	"github.com/fatih/structs"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupLocalRoutes(s service.Service, app *fiber.App) {
	app.Get("/ping", pong)
	app.Get("/", GetMain(s.DB))
	app.Get("/keys", GetKeys)
	app.Post("/keys", PostKeys)
}

// Healthcheck
func pong(c *fiber.Ctx) error {
	return c.SendString("pong")
}

func GetMain(db *gorm.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		localData := GetLocalMessages(db)
		return c.Render("main/root", structs.Map(localData))

	}
}
