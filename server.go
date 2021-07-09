package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/ewenquim/horkruxes/api"
	"github.com/ewenquim/horkruxes/service"
	"github.com/ewenquim/horkruxes/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/template/html"
)

func runServer() {
	fsub, _ := fs.Sub(staticFS, "static") // error ignored because it can only happen if binary is not correctly built

	serverConfig := loadServerConfig()

	// Database setup
	db := initDatabase()
	sqldb, err := db.DB()
	if err != nil {
		panic("the db have an issue") // should not happen with sqlite, might happen with server databases (MySQL, Pg...)
	}
	defer sqldb.Close()

	// Service init
	s := service.Service{
		DB:           db,
		ServerConfig: serverConfig,
		Regexes:      service.InitializeDetectors(),
	}

	// Templating engine init
	engine := html.NewFileSystem(http.FS(templatesFS), ".html")
	engine.Debug(s.ServerConfig.Debug)

	engine.AddFunc("md", service.MarkDowner)

	// Server and middlewares
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Limit API posts to 5/day (still can use local post without limitations)
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return (c.Method() == "GET" || c.Path() != "/api/message")
		},
		Max:        5,
		Expiration: 24 * time.Hour,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "everyone" // does not depend on c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
	}))

	app.Use(helmet.New())
	app.Use(cors.New())
	// app.Use(csrf.New()) // Useless and blocks post requests...

	// app.Use(favicon.New(favicon.Config{FileSystem: http.FS(fsub)}))

	if s.ServerConfig.Debug {
		app.Use(logger.New())
	}

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
		return c.Status(fiber.StatusNotFound).SendString("Sorry, can't find that! Check your URL")
	})

	app.Listen(fmt.Sprintf(":%v", s.ServerConfig.Port))

}