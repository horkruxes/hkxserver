package views

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fatih/structs"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/horkruxes/hkxserver/client"
	"github.com/horkruxes/hkxserver/model"
	"github.com/horkruxes/hkxserver/query"
	"github.com/horkruxes/hkxserver/service"
)

func GetKeys(c *fiber.Ctx) error {
	outputData := client.GenKeys()
	return c.Render("templates/keys/root", structs.Map(outputData))
}

func PostKeys(c *fiber.Ctx) error {
	outputData := client.GenKeys()

	// Get form data and reinject into output data
	outputData.Sig = strings.TrimSpace(c.FormValue("signature"))
	outputData.Sec = strings.TrimSpace(c.FormValue("secret-key"))
	outputData.Pub = strings.TrimSpace(c.FormValue("public-key"))
	outputData.DisplayedName = strings.TrimSpace(c.FormValue("displayed-name"))
	outputData.Content = strings.TrimSpace(c.FormValue("message"))
	outputData.MessageID = strings.TrimSpace(c.FormValue("answer-to"))
	outputData.Verif = true

	if outputData.Sig == "" {
		// Answers to the signature GENERATION form
		outputData.Sig = client.SignStrings(outputData.Sec, outputData.Pub, outputData.DisplayedName, outputData.Content, outputData.MessageID)
		outputData.Valid = client.VerifyFromString(outputData.Pub, outputData.Sig, outputData.DisplayedName, outputData.Content, outputData.MessageID)
		outputData.Sec = ""
	} else {
		// Answers to the signature VERIFICATION form
		outputData.Valid = client.VerifyFromString(outputData.Pub, outputData.Sig, outputData.DisplayedName, outputData.Content, outputData.MessageID)
		outputData.Sig = ""
	}

	return c.Render("templates/keys/root", structs.Map(outputData))
}

func GetFaq(c *fiber.Ctx) error {
	return c.Render("templates/faq/root", fiber.Map{})
}

func GetMainNoFront(c *fiber.Ctx) error {
	return c.SendString(`The pod admin chose to only use the 'data' part of Horkruxes.
Sorry, you'll have to use another client to see the messages.`)
}

func GetMain(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		s.ClientConfig = parseFormsToService(c, s)
		localData := GetMessagesAndMainPageInfo(s)
		return c.Render("templates/main/root", structs.Map(localData))
	}
}

func GetAuthor(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		s.ClientConfig = parseFormsToService(c, s)
		pubKey := c.Params("pubKey")
		fmt.Println("pods list", s.GeneralConfig.TrustedPods)
		localData := GetAuthorMessagesAndMainPageInfo(s, pubKey)
		return c.Render("templates/main/root", structs.Map(localData))
	}
}

func GetComments(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		s.ClientConfig = parseFormsToService(c, s)
		id := c.Params("uuid")
		_, err := uuid.Parse(id)
		if err != nil {
			return c.Status(400).SendString("Invalid UUID" + err.Error())
		}
		localData := GetCommentsAndMainPageInfo(s, id)
		return c.Render("templates/main/root", structs.Map(localData))
	}
}

func NewMessage(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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

		// Sign if necessary
		if secretKey := strings.TrimSpace(c.FormValue("secret-key")); secretKey != "" && payload.SignatureBase64 == "" {
			fmt.Println("signing server-side")
			payload.SignatureBase64, err = client.SignMessage(payload, secretKey)
			if err != nil {
				fmt.Println("can't sign payload:", err)
				return c.Status(500).SendString(err.Error())
			}
			fmt.Println("signed:", payload.SignatureBase64)
		}

		newMessage := model.Message{}

		fmt.Println("try to post to:", payload.Pod)
		// Check if can do the db operations right now or if it should transfer the payload to another API
		if payload.Pod == "" {

			newMessage, err = query.NewMessage(s, payload)
			if err != nil {
				fmt.Printf("error uploading message locally: %v\n%+v\n", err, payload)
				return c.Status(422).SendString(err.Error())
			}
		} else {
			reader, err := json.Marshal(payload)
			if err != nil {
				fmt.Println("err:", err)
			}
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

				err = json.Unmarshal(body, &newMessage)
				if err != nil {
					fmt.Println("err:", err)
				}

				return c.SendString("error sending the message " + string(body))
			}
		}

		if newMessage.MessageID == "" {
			return c.Redirect("/comments/" + newMessage.ID)
		} else {
			return c.Redirect("/comments/" + newMessage.MessageID)
		}
	}
}
