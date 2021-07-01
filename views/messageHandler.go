package views

import (
	"encoding/base64"
	"fmt"

	"github.com/ewenquim/horkruxes-client/model"
	"github.com/fatih/structs"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetMessages(db *gorm.DB) model.PageData {
	// call := []string{}

	messages := model.GetMessagesFromDB(db)

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
	return data
}

func GetKeys(c *fiber.Ctx) error {
	outputData := model.GenKeys()
	return c.Render("keys/root", structs.Map(outputData))
}

func PostKeys(c *fiber.Ctx) error {

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
}
