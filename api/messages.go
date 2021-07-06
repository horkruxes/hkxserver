package api

import (
	"encoding/base64"
	"strings"

	"github.com/ewenquim/horkruxes/exceptions"
	"github.com/ewenquim/horkruxes/model"
	"github.com/ewenquim/horkruxes/service"
	"github.com/gofiber/fiber/v2"
)

type NewMessagePayload struct {
	Signature string
	PublicKey string `json:"public-key"`
	Content   string
	Name      string
	Pod       string
	MessageID string `json:"messageID,omitempty"`
}

func GetMessagesJSON(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		data := model.GetMessagesFromDB(s)
		return c.Status(201).JSON(data)
	}
}

func GetCommentsJSON(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		data := model.GetCommentsTo(s, id)
		return c.Status(201).JSON(data)
	}
}

func GetMessageJSON(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		message, err := model.GetMessageFromDB(s, id)
		if err != nil {
			return err
		}
		return c.JSON(message)
	}
}

func GetMessagesFromAuthorJSON(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("pubKey")
		pubKey := service.SafeURLToBase64(id)
		message := model.GetMessagesFromAuthor(s, pubKey)
		return c.JSON(message)
	}
}

func NewMessage(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Get message
		payload := NewMessagePayload{}

		// Read Body
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		// Translate into Message struct and verify conditions
		message, err := PayloadToValidMessage(payload)
		if err != nil {
			return c.Status(409).SendString(err.Error())
		}

		// Register
		model.NewMessage(s, message)
		return c.Redirect("/")
	}
}

func PayloadToValidMessage(payload NewMessagePayload) (model.Message, error) {
	message := model.Message{}

	var err error

	message.SignatureBase64 = strings.TrimSpace(payload.Signature)
	message.Signature, err = base64.StdEncoding.DecodeString(message.SignatureBase64)
	if err != nil {
		return message, exceptions.WrongSignature
	}
	message.AuthorBase64 = strings.TrimSpace(payload.PublicKey)
	message.AuthorPubKey, err = base64.StdEncoding.DecodeString(message.AuthorBase64)
	if err != nil {
		return message, exceptions.WrongSignature
	}
	message.Content = strings.TrimSpace(payload.Content)
	message.DisplayedName = strings.TrimSpace(payload.Name)
	message.MessageID = strings.TrimSpace(payload.MessageID)

	if !message.VerifyConditions() {
		return message, exceptions.ErrorRecordTooLongFound
	}
	if !message.VerifyOwnerShip() {
		return message, exceptions.WrongSignature

	}
	message.Correct = true
	return message, nil
}
