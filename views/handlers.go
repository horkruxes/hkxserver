package views

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ewenquim/horkruxes/api"
	"github.com/ewenquim/horkruxes/model"
	"github.com/ewenquim/horkruxes/service"
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

func GetMain(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		localData := GetMessagesAndMainPageInfo(s)
		return c.Render("main/root", structs.Map(localData))
	}
}

func GetAuthor(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := service.SafeURLToBase64(c.Params("pubKey"))
		localData := GetAuthorMessagesAndMainPageInfo(s, id)
		return c.Render("main/root", structs.Map(localData))
	}
}

func GetComments(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("uuid")
		localData := GetCommentsAndMainPageInfo(s, id)
		return c.Render("main/root", structs.Map(localData))
	}
}

func NewMessage(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		payload := FromFormToPayload(c)

		reader, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("err:", err)
		}

		fmt.Println("try to post to:", payload.Pod)
		// Check if can do the db operations right now or if it should transfer the payload to another API
		if payload.Pod == "" {
			message, statusCode, err := api.PayloadToValidMessage(s, payload)
			if err != nil {
				return c.Status(statusCode).SendString(err.Error())
			}
			fmt.Println("new msg", message)
			err = model.NewMessage(s, message)
			if err != nil {
				return c.Status(statusCode).SendString(err.Error())
			}
		} else {
			resp, err := http.Post("https://"+payload.Pod+"/api/message", "application/json", bytes.NewBuffer(reader))
			if err != nil {
				fmt.Println("err:", err)
			}
			if resp.StatusCode == 409 {

				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("err:", err)
				}
				fmt.Println(string(body))

				c.SendString("error sending the message " + string(body))
			}
		}

		return c.Redirect("/")
	}
}
