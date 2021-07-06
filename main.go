package main

import (
	"fmt"

	"github.com/ewenquim/horkruxes/api"
	"github.com/ewenquim/horkruxes/service"
	"github.com/ewenquim/horkruxes/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/redirect/v2"
	"github.com/gofiber/template/html"
)

type SignMessage struct {
	AuthorPubKey string // base64-encoded
	AuthorSecKey string // base64-encoded
	Content      string
	Signature    string
}

func main() {
	// Database setup
	db := initDatabase()
	sqldb, _ := db.DB()
	defer sqldb.Close()

	s := service.Service{
		DB:           initDatabase(),
		ServerConfig: loadServerConfig(),
	}

	// Server and middlewares
	engine := html.New("./templates", ".html")
	engine.AddFunc("md", service.MarkDowner)
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(cors.New())

	app.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"/messages": "/",
		},
		StatusCode: 301,
	}))

	app.Use(logger.New())

	// Static routes
	app.Static("", "./static")
	fmt.Println("Static server started")

	// Backend - DB operations routes (potentially online)
	api.SetupApiRoutes(s, app)
	fmt.Println("API started")

	// Frontend - Local views and template rendering
	views.SetupLocalRoutes(s, app)
	fmt.Println("Frontend started")

	// 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Sorry, can't find that! Check your URL")
	})

	app.Listen(fmt.Sprintf(":%v", s.ServerConfig.Port))
}
