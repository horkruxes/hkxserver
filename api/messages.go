package api

import (
	"encoding/base64"

	"github.com/ewenquim/horkruxes-client/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetMessagesJSON(db *gorm.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		data := model.GetMessagesFromDB(db)
		c.Status(201).JSON(fiber.Map{"response": data})
		return nil
	}
}

func GetMessageJSON(db *gorm.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		message := model.GetMessageFromDB(db, id)
		return c.JSON(message)
	}
}

func NewMessage(db *gorm.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		message := &model.Message{}
		var err error
		message.SignatureBase64 = c.FormValue("signature")
		message.Signature, err = base64.StdEncoding.DecodeString(message.SignatureBase64)
		if err != nil {
			return c.Status(503).SendString("error wrong signature")
		}
		message.AuthorBase64 = c.FormValue("public-key")
		message.AuthorPubKey, err = base64.StdEncoding.DecodeString(message.AuthorBase64)
		if err != nil {
			return c.Status(503).SendString("error wrong public key")
		}
		message.Content = c.FormValue("message")
		message.DisplayedName = c.FormValue("name")

		if message.VerifyOwnerShip() {
			message.Correct = true
			model.NewMessage(db, message)
			return c.Redirect("/")
		}
		return c.Status(503).SendString("error unvalid public key/signature")
	}
}
