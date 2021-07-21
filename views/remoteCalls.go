package views

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ewenquim/horkruxes/model"
	"github.com/ewenquim/horkruxes/service"
)

func getMessagesFrom(s service.Service, path string) []model.Message {
	var messages []model.Message
	for _, pod := range s.ClientConfig.AllPodsList {
		if pod.URL != s.GeneralConfig.URL {
			adress := "https://" + pod.URL + path

			fmt.Println("rq to", adress)
			resp, err := http.Get(adress)
			if err != nil {
				fmt.Println("err", err)
				continue
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("err", err)
				continue
			}
			if !json.Valid(body) {
				continue
			}

			msg := &[]model.Message{}
			err = json.Unmarshal(body, msg)
			if err != nil {
				fmt.Println("err", err)
				continue
			}
			messages = append(messages, *msg...)
		}
	}
	return messages
}

func getSingleMessageFromEachPod(s service.Service, path string) []model.Message {
	var messages []model.Message
	for _, pod := range s.ClientConfig.AllPodsList {
		if pod.URL != s.GeneralConfig.URL {
			adress := "https://" + pod.URL + path

			fmt.Println("rq to", adress)
			resp, err := http.Get(adress)
			if err != nil {
				fmt.Println("err", err)
				continue
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("err", err)
				continue
			}
			if !json.Valid(body) {
				continue
			}

			msg := &model.Message{}
			err = json.Unmarshal(body, msg)
			if err != nil {
				fmt.Println("err", err)
				continue
			}
			messages = append(messages, *msg)
		}
	}
	return messages
}
