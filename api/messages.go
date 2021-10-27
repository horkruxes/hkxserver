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

// GetAllJSON godoc
// @Id all
// @Summary Show all messages (original messages and comments) without any more structure
// @Description get string by ID
// @Produce application/json
// @Success 200 {array} model.Message
// @Router /all [get]
func GetAllJSON(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		data := model.GetAllFromDB(s)
		return c.Status(200).JSON(data)
	}
}

// GetMessagesJSON godoc
// @Id messages
// @Summary Show all original messages (for the front page)
// @Description get string by ID
// @Produce application/json
// @Success 200 {array} model.Message
// @Router /message [get]
func GetMessagesJSON(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		data := model.GetMessagesFromDB(s)
		return c.Status(200).JSON(data)
	}
}

// GetCommentsJSON godoc
// @Summary Show comments to a specific message
// @Description get string by ID
// @Produce application/json
// @Param uuid path string true "UUID of original message corresponding to comments"
// @Success 200 {array} model.Message
// @Router /comments/{uuid} [get]
func GetCommentsJSON(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		data := model.GetCommentsTo(s, id)
		return c.Status(200).JSON(data)
	}
}

// GetCommentsJSON godoc
// @Summary Show a specific message
// @Description get string by ID
// @Produce application/json
// @Param uuid path string true "Message UUID"
// @Success 200 {object} model.Message
// @Router /message/{uuid} [get]
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

// GetMessagesFromAuthorJSON godoc
// @Summary Show messages of a specific user
// @Description get string by ID
// @Produce application/json
// @Param pubkey path string true "Author ed25519 public key"
// @Success 200 {array} model.Message
// @Router /user/{pubkey} [get]
func GetMessagesFromAuthorJSON(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("pubKey")
		pubKey := service.SafeURLToBase64(id)
		message := model.GetMessagesFromAuthor(s, pubKey)
		return c.JSON(message)
	}
}

// NewMessage godoc
// @Summary Post a new message
// @Description get string by ID
// @Produce json
// @Router /new [post]
func NewMessage(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		fmt.Println("Received POST request to create new message")
		// Get message
		payload := model.Message{}

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

func PayloadToValidMessage(s service.Service, payload model.Message) (model.Message, int, error) {
	message := model.Message{}

	var err error

	message.SignatureBase64 = strings.TrimSpace(payload.SignatureBase64)
	_, err = base64.StdEncoding.DecodeString(message.SignatureBase64)
	if err != nil {
		return message, fiber.StatusBadRequest, exceptions.ErrorWrongSignature
	}
	message.AuthorBase64 = strings.TrimSpace(payload.AuthorBase64)
	_, err = base64.StdEncoding.DecodeString(message.AuthorBase64)
	if err != nil {
		return message, fiber.StatusBadRequest, exceptions.ErrorWrongSignature
	}
	message.Content = strings.TrimSpace(payload.Content)
	message.DisplayedName = strings.TrimSpace(payload.DisplayedName)
	message.MessageID = strings.TrimSpace(payload.MessageID)

	if statusCode, err := message.VerifyConditions(s); err != nil {
		return message, statusCode, err
	}
	return message, fiber.StatusOK, nil
}
