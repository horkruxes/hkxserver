package api

import (
	"encoding/base64"
	"strings"

	"github.com/ewenquim/horkruxes/exceptions"
	"github.com/ewenquim/horkruxes/model"
	"github.com/ewenquim/horkruxes/service"
	"github.com/ewenquim/horkruxes/views"
	"github.com/gofiber/fiber/v2"
)

func GetMessagesJSON(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		data := model.GetMessagesFromDB(s)
		return c.Status(201).JSON(data)
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
		pubKey := views.SafeURLToBase64(id)
		message := model.GetMessagesFromAuthor(s, pubKey)
		return c.JSON(message)
	}
}

func NewMessage(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		message := &model.Message{}
		var err error

		message.SignatureBase64 = strings.TrimSpace(c.FormValue("signature"))
		message.Signature, err = base64.StdEncoding.DecodeString(message.SignatureBase64)
		if err != nil {
			return c.Status(409).SendString("error wrong signature")
		}
		message.AuthorBase64 = strings.TrimSpace(c.FormValue("public-key"))
		message.AuthorPubKey, err = base64.StdEncoding.DecodeString(message.AuthorBase64)
		if err != nil {
			return c.Status(409).SendString("error wrong public key")
		}
		message.Content = strings.TrimSpace(c.FormValue("message"))
		message.DisplayedName = strings.TrimSpace(c.FormValue("name"))

		if !message.VerifyConditions() {
			return c.Status(409).SendString(exceptions.ErrorRecordTooLongFound.Error())
		}
		if !message.VerifyOwnerShip() {
			return c.Status(409).SendString("error unvalid public key/signature")

		}
		message.Correct = true
		model.NewMessage(s, message)
		return c.Redirect("/")
	}
}
