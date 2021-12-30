package api

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/horkruxes/hkxserver/client"
	"github.com/horkruxes/hkxserver/exceptions"
	"github.com/horkruxes/hkxserver/model"
	"github.com/horkruxes/hkxserver/query"
	"github.com/horkruxes/hkxserver/service"
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
		data := query.GetAll(s)
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
		data := query.GetMessages(s)
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
		data := query.GetCommentsTo(s, id)
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
		message, err := query.GetMessage(s, id)
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
		pubKey := c.Params("pubKey")
		message := query.GetMessagesFromAuthor(s, pubKey)
		return c.JSON(message)
	}
}

// NewMessage godoc
// @Summary Post a new message
// @Description Post a new message
// @Produce json
// @Router /message [post]
func NewMessage(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		fmt.Println("Received POST request to create new message")
		// Get message
		payload := model.Message{}

		// Read Body
		if err := c.BodyParser(&payload); err != nil {
			fmt.Println("can't parse payload:", err)
			return c.Status(500).SendString(err.Error())
		}
		fmt.Println("content:", payload)
		secretKey := strings.TrimSpace(c.FormValue("secret-key"))
		if secretKey != "" && strings.TrimSpace(c.FormValue("signature")) == "" {
			payload.SignatureBase64 = client.SignMessage(secretKey, payload.AuthorBase64, payload.DisplayedName, payload.Content, payload.MessageID)
		}

		// Translate into Message struct and verify conditions
		message, statusCode, err := PayloadToValidMessage(s, payload)
		if err != nil {
			fmt.Println("err:", err)
			return c.Status(statusCode).SendString(err.Error())
		}

		// Register
		err = query.NewMessage(s, message)
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
	_, err = base64.URLEncoding.DecodeString(message.SignatureBase64)
	if err != nil {
		return message, 400, exceptions.ErrorWrongSignature
	}
	message.AuthorBase64 = strings.TrimSpace(payload.AuthorBase64)
	_, err = base64.URLEncoding.DecodeString(message.AuthorBase64)
	if err != nil {
		return message, 400, exceptions.ErrorWrongSignature
	}
	message.Content = strings.TrimSpace(payload.Content)
	message.DisplayedName = strings.TrimSpace(payload.DisplayedName)
	message.MessageID = strings.TrimSpace(payload.MessageID)

	if err := message.VerifyConstraints(); err != nil {
		return message, 400, err
	}
	if err := query.VerifyServerConstraints(s, message); err != nil {
		return message, 400, err
	}
	return message, fiber.StatusOK, nil
}
