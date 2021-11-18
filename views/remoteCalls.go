package views

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/horkruxes/hkxserver/model"
	"github.com/horkruxes/hkxserver/service"
)

func getMessagesFrom(s service.Service, path string) []model.Message {
	var messages []model.Message
	for _, pod := range s.ClientConfig.AllPodsList {
		if pod.URL != s.GeneralConfig.URL {
			adress := "https://" + pod.URL + path

			fmt.Println("rq to", adress)
			//#nosec TODO avoid SSRF
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
			//#nosec TODO avoid SSRF
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
