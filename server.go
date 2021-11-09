package main

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/template/html"
	"github.com/horkruxes/hkxserver/api"
	"github.com/horkruxes/hkxserver/service"
	"github.com/horkruxes/hkxserver/views"

	// Swagger middlewre and docs generated by `swag init`
	swagger "github.com/arsmn/fiber-swagger/v2"
	_ "github.com/horkruxes/hkxserver/docs"
)

func setupService() service.Service {
	// Database setup
	db := initDatabase(dbOptions{})

	// Service init
	s := loadConfig()
	s.DB = db
	s.Regexes = service.InitializeDetectors()
	return s
}

func setupServer(s service.Service) (fiber.App, int64) {
	fsub, _ := fs.Sub(staticFS, "static") // error ignored because it can only happen if binary is not correctly built

	// Templating engine init
	engine := html.NewFileSystem(http.FS(templatesFS), ".html")
	engine.Debug(s.ServerConfig.Debug)

	engine.AddFunc("md", service.MarkDowner)

	// Server and middlewares
	app := fiber.New(fiber.Config{
		Views:   engine,
		AppName: "Horkruxes",
	})

	app.Use(helmet.New())
	app.Use(cors.New())
	// app.Use(csrf.New()) // Useless and blocks post requests...

	// app.Use(favicon.New(favicon.Config{FileSystem: http.FS(fsub)}))

	if s.ServerConfig.Debug {
		app.Use(logger.New())
	} else {
		app.Use(cache.New())
	}

	// Swagger
	app.Get("/swagger/*", swagger.Handler)
	app.Get("/swagger", func(c *fiber.Ctx) error { return c.Redirect("/swagger/") })

	// Backend - DB operations routes (potentially online)
	api.SetupApiRoutes(s, app)
	fmt.Println("API started")

	// Frontend - Local views and template rendering
	views.SetupLocalRoutes(s, app)
	fmt.Println("Frontend started")

	// Static routes
	app.Use(filesystem.New(filesystem.Config{
		Root: http.FS(fsub),
	}))
	// app.Static("", "./static")
	fmt.Println("Static server started")

	// 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("404 error: wrong URL")
	})

	return *app, s.ServerConfig.Port
}
