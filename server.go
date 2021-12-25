package main

//go:generate npx tailwindcss -i tailwind.css -o static/tailwindstyles.css -m

import (
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/template/html"
	"github.com/horkruxes/hkxserver/api"
	"github.com/horkruxes/hkxserver/service"
	"github.com/horkruxes/hkxserver/views"
	"github.com/microcosm-cc/bluemonday"

	// Swagger middlewre and docs generated by `swag init`
	swagger "github.com/arsmn/fiber-swagger/v2"
	_ "github.com/horkruxes/hkxserver/docs"
)

func setupService() service.Service {
	// Database setup
	db := initDatabase("db.sqlite3")

	// Service init
	s := loadConfig()
	s.DB = db
	s.Regexes = service.InitializeDetectors()
	s.ContentPolicy = bluemonday.UGCPolicy()
	return s
}

func setupServer(s service.Service) (fiber.App, int64) {
	fsub, _ := fs.Sub(staticFS, "static") // error ignored because it can only happen if binary is not correctly built

	// Templating engine init
	engine := html.NewFileSystem(http.FS(templatesFS), ".html")
	engine.Debug(s.ServerConfig.Debug)

	engine.AddFunc("md", views.MarkDowner(s.ContentPolicy))

	// Server and middlewares
	app := fiber.New(fiber.Config{
		Views:                   engine,
		AppName:                 "Horkruxes",
		ServerHeader:            "hkxserver",
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"127.0.0.1", "0.0.0.0"},
		ProxyHeader:             "X-Forwarded-For",
	})

	// Security
	app.Use(helmet.New(helmet.Config{
		HSTSPreloadEnabled:    true,
		HSTSMaxAge:            31536000,
		HSTSExcludeSubdomains: true,
		ReferrerPolicy:        "same-origin",
		ContentSecurityPolicy: "default-src 'self'; img-src https: data:;",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:               50,
		Expiration:        1 * time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	app.Use(cors.New())

	if s.ServerConfig.Debug {
		app.Use(logger.New())
	} else {
		app.Use(cache.New())
	}

	app.Use(compress.New())

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
		Root:   http.FS(fsub),
		MaxAge: 604_800, // 1 week cache (7*24*3600 = 604 800 seconds)
	}))
	// app.Static("", "./static")
	fmt.Println("Static server started")

	// 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("404 error: wrong URL")
	})

	return *app, s.ServerConfig.Port
}
