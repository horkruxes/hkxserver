package model

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	"github.com/fatih/structs"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	AuthorPubKey    []byte `json:"authorPubKey"`
	AuthorBase64    string `json:"authorBase64"`
	DisplayedName   string `json:"displayedName"` // Name Chosen by author, bo restriction
	Content         string `json:"content"`
	Signature       []byte `json:"signature"`
	SignatureBase64 string `json:"signatureBase64"`
	Correct         bool
	Color           string
	// Date          time.Time
}

type PageData struct {
	PageTitle string
	Messages  []Message
}

type Service struct {
	DB *gorm.DB
}

func (s Service) SetupRoutes(app *fiber.App) {
	app.Get("/test", s.GetMessages)
	app.Get("/api/message/:id", s.GetMessage)
	app.Post("/api/message", s.NewMessage)
	app.Delete("/api/message/:id", s.DeleteMessage)
}

func (s Service) GetMessages(c *fiber.Ctx) error {
	var messages []Message
	s.DB.Where("correct = ?", true).Find(&messages)
	data := PageData{Messages: messages}
	for i, message := range data.Messages {
		data.Messages[i].Correct = message.verifyOwnerShip()
		fmt.Println(data.Messages[i].Correct)
		data.Messages[i].AuthorBase64 = base64.StdEncoding.EncodeToString(message.AuthorPubKey)
		data.Messages[i].Color = ColorFromString(string(message.AuthorPubKey))
		data.Messages[i].SignatureBase64 = base64.StdEncoding.EncodeToString(message.Signature)
	}

	return c.Render("main/root", structs.Map(data))
}

func (s Service) GetMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	var message Message
	s.DB.Find(&message, id)
	return c.JSON(message)
}

func (s Service) NewMessage(c *fiber.Ctx) error {
	message := &Message{}

	if err := c.BodyParser(message); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	if message.verifyOwnerShip() {
		message.Correct = true
		s.DB.Create(&message)
		return c.JSON(message)
	}
	return c.Status(503).SendString("error wrong signature")
}

func (s Service) DeleteMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	var message Message
	s.DB.First(&message, id)
	if string(message.AuthorPubKey) == "" {
		return c.Status(500).SendString("No Book Found with ID")
	}
	s.DB.Delete(&message)
	return c.SendString("Book Successfully deleted")
}

func (message Message) verifyOwnerShip() bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	fmt.Println("pub", message.AuthorPubKey)
	fmt.Println("sig", message.Signature)

	if message.AuthorPubKey == nil || message.Signature == nil || len(message.AuthorPubKey) == 0 || len(message.Signature) == 0 {
		return false
	}
	return ed25519.Verify(message.AuthorPubKey, []byte(message.Content+string(message.AuthorPubKey)), message.Signature)
}
