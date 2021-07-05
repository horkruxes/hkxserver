package views

import (
	"strings"

	"github.com/ewenquim/horkruxes/model"
	"github.com/fatih/structs"
	"github.com/gofiber/fiber/v2"
)

func GetKeys(c *fiber.Ctx) error {
	outputData := model.GenKeys()
	return c.Render("keys/root", structs.Map(outputData))
}

func PostKeys(c *fiber.Ctx) error {

	outputData := model.GenKeys()

	// Get form data and reinject into output data
	outputData.Sig = strings.TrimSpace(c.FormValue("signature"))
	outputData.Sec = strings.TrimSpace(c.FormValue("secret-key"))
	outputData.Pub = strings.TrimSpace(c.FormValue("public-key"))
	outputData.Content = strings.TrimSpace(c.FormValue("message"))
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

func GetFaq(c *fiber.Ctx) error {
	return c.Render("faq/root", fiber.Map{})
}
