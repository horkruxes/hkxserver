package api

import (
	"encoding/base64"

	"github.com/ewenquim/horkruxes/exceptions"
	"github.com/ewenquim/horkruxes/model"
	"github.com/ewenquim/horkruxes/service"
	"github.com/gofiber/fiber/v2"
)

func GetMessagesJSON(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		data := model.GetMessagesFromDB(s)
		c.Status(201).JSON(fiber.Map{"response": data})
		return nil
	}
}

func GetMessageJSON(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		message := model.GetMessageFromDB(s, id)
		return c.JSON(message)
	}
}

func GetMessagesFromAuthorJSON(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("pubKey")
		message := model.GetMessagesFromAuthor(s, id)
		return c.JSON(message)
	}
}

func NewMessage(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		message := &model.Message{}
		var err error
		message.SignatureBase64 = c.FormValue("signature")
		message.Signature, err = base64.StdEncoding.DecodeString(message.SignatureBase64)
		if err != nil {
			return c.Status(409).SendString("error wrong signature")
		}
		message.AuthorBase64 = c.FormValue("public-key")
		message.AuthorPubKey, err = base64.StdEncoding.DecodeString(message.AuthorBase64)
		if err != nil {
			return c.Status(409).SendString("error wrong public key")
		}
		message.Content = c.FormValue("message")
		message.DisplayedName = c.FormValue("name")

		if !message.VerifyConditions() {
			return c.Status(409).SendString(exceptions.ErrorRecordTooLongFound.Error())
		}
		if message.VerifyOwnerShip() {
			message.Correct = true
			model.NewMessage(s, message)
			return c.Redirect("/")
		}
		return c.Status(409).SendString("error unvalid public key/signature")
	}
}
