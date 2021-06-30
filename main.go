package main

import (
	"encoding/base64"
	"fmt"

	"github.com/ewenquim/horkruxes-client/api"
	"github.com/ewenquim/horkruxes-client/model"
	"github.com/fatih/structs"
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
	// data := mock()

	// Server and middlewares
	engine := html.New("./templates", ".html")
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

	// Database setup
	db := initDatabase()
	sqldb, _ := db.DB()
	defer sqldb.Close()

	service := model.Service{
		DB: initDatabase(),
	}

	// DB operations routes
	api.SetupMessagesRoutes(service, app)

	app.Static("/static", "./static")

	// Healthcheck
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	app.Get("/", func(c *fiber.Ctx) error {
		// call := []string{}

		messages := service.GetMessagesFromDB()

		// for i, ip := range call {

		// 	messages = append(messages, )
		// }

		// s.DB.Where("correct = ?", true).Find(&messages)
		fmt.Println(messages)
		data := model.PageData{Messages: messages}
		for i, message := range data.Messages {
			data.Messages[i].Correct = message.VerifyOwnerShip()
			fmt.Println(data.Messages[i].Correct)
			data.Messages[i].AuthorBase64 = base64.StdEncoding.EncodeToString(message.AuthorPubKey)
			data.Messages[i].Color = model.ColorFromString(string(message.AuthorPubKey))
			data.Messages[i].SignatureBase64 = base64.StdEncoding.EncodeToString(message.Signature)
		}

		return c.Render("main/root", structs.Map(data))
	})

	app.Get("/keys", func(c *fiber.Ctx) error {
		outputData := model.GenKeys()
		return c.Render("keys/root", structs.Map(outputData))
	})

	app.Post("/keys", func(c *fiber.Ctx) error {
		outputData := model.GenKeys()

		// Get form data and reinject into output data
		outputData.Sig = c.FormValue("signature")
		outputData.Sec = c.FormValue("secret-key")
		outputData.Pub = c.FormValue("public-key")
		outputData.Content = c.FormValue("message")
		outputData.Verif = true

		if outputData.Sig == "" {
			// Answers to the signature GENERATION form
			outputData.Sig = model.SignMessage(outputData.Sec, outputData.Pub, outputData.Content)
			outputData.Valid = model.VerifyFromString(outputData.Pub, outputData.Sig, outputData.Content)
			outputData.Sec = ""
		} else {
			// Answers to the signature VERIFICATION form
			outputData.Valid = model.VerifyFromString(outputData.Pub, outputData.Sig, outputData.Content)
			outputData.Sig = ""
		}

		return c.Render("keys/root", structs.Map(outputData))
	})

	// 404
	app.Use("404", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Sorry can't find that!")
	})

	app.Listen(":80")

}
