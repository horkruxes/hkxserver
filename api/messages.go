package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
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
// @Param message body model.Message true "Message"
// @Router /message [post]
func NewMessage(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		fmt.Println("Received POST request to create new message")

		// Read Body
		payload := model.Message{}
		if err := c.BodyParser(&payload); err != nil {
			fmt.Println("can't parse payload:", err)
			return c.Status(500).SendString(err.Error())
		}
		fmt.Println("content:", payload)

		// Normalize
		err := payload.Normalize()
		if err != nil {
			fmt.Println("can't normalize payload:", err)
			return c.Status(500).SendString(err.Error())
		}

		// Register message
		newMessage, err := query.NewMessage(s, payload)
		if err != nil {
			fmt.Println("error:", err)
			return c.Status(400).SendString(err.Error())
		}
		return c.JSON(newMessage)
	}
}
