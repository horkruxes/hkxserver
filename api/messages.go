package api

import (
	"encoding/base64"
	"fmt"
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
		return c.Status(200).JSON(data)
	}
}

func GetCommentsJSON(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		data := model.GetCommentsTo(s, id)
		return c.Status(200).JSON(data)
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
		fmt.Println("Received POST request to create new message")
		// Get message
		payload := NewMessagePayload{}

		// Read Body
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		fmt.Println("content:", payload)

		// Translate into Message struct and verify conditions
		message, statusCode, err := PayloadToValidMessage(s, payload)
		if err != nil {
			fmt.Println("err:", err)
			return c.Status(statusCode).SendString(err.Error())
		}

		// Register
		err = model.NewMessage(s, message)
		if err != nil {
			fmt.Println("err:", err)
			return c.Status(statusCode).SendString(err.Error())
		}
		return c.Redirect("/")
	}
}

func PayloadToValidMessage(s service.Service, payload NewMessagePayload) (model.Message, int, error) {
	message := model.Message{}

	var err error

	message.SignatureBase64 = strings.TrimSpace(payload.Signature)
	message.Signature, err = base64.StdEncoding.DecodeString(message.SignatureBase64)
	if err != nil {
		return message, fiber.StatusBadRequest, exceptions.ErrorWrongSignature
	}
	message.AuthorBase64 = strings.TrimSpace(payload.PublicKey)
	message.AuthorPubKey, err = base64.StdEncoding.DecodeString(message.AuthorBase64)
	if err != nil {
		return message, fiber.StatusBadRequest, exceptions.ErrorWrongSignature
	}
	message.Content = strings.TrimSpace(payload.Content)
	message.DisplayedName = strings.TrimSpace(payload.Name)
	message.MessageID = strings.TrimSpace(payload.MessageID)

	if statusCode, err := message.VerifyConditions(s); err != nil {
		return message, statusCode, err
	}
	message.Correct = true
	return message, fiber.StatusOK, nil
}
