package views

import (
	"encoding/base64"

	"github.com/ewenquim/horkruxes-client/model"
	"github.com/ewenquim/horkruxes-client/service"
	"github.com/fatih/structs"
	"github.com/gofiber/fiber/v2"
)

type PageData struct {
	Server   ServerData
	Messages []model.Message
}

type ServerData struct {
	Name string
	IP   string
}

// Get Local and online messages, checks validity and return view
func GetMessagesAndMainPageInfo(s service.Service) PageData {

	// Get local messages
	messages := model.GetMessagesFromDB(s.DB)

	// Get other pods messages
	// call := []string{}
	// for i, ip := range call {
	// 	messages = append(messages, )
	// }

	// Check data validity
	for i, message := range messages {
		messages[i].Correct = message.VerifyOwnerShip()
		// fmt.Println(messages[i].Correct)
		messages[i].AuthorBase64 = base64.StdEncoding.EncodeToString(message.AuthorPubKey)
		messages[i].Color = model.ColorFromString(string(message.AuthorPubKey))
		messages[i].SignatureBase64 = base64.StdEncoding.EncodeToString(message.Signature)
	}

	// Inject view
	return PageData{
		Messages: messages,
		Server:   ServerData{Name: s.ServerConfig.Name, IP: s.ServerConfig.URL},
	}
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
