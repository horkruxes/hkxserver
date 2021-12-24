package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/horkruxes/hkxserver/model"
)

func GetMessagesFrom(servers []string, path string) []model.Message {
	chanMessages := make(chan []model.Message, len(servers))

	var wg sync.WaitGroup
	for _, pod := range servers {
		wg.Add(1)

		go func(chanMessages chan<- []model.Message, pod string) {
			defer wg.Done()

			adress := "https://" + pod + path

			fmt.Println("rq to", adress)
			//#nosec
			resp, err := http.Get(adress)
			if err != nil {
				fmt.Println("err", err)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("err", err)
				return
			}
			if !json.Valid(body) {
				return
			}

			msg := []model.Message{}
			err = json.Unmarshal(body, &msg)
			if err != nil {
				fmt.Println("err", err)
				return
			}
			chanMessages <- msg

		}(chanMessages, pod)
	}

	var messages []model.Message

	wg.Wait()
	close(chanMessages)
	for msgs := range chanMessages {
		messages = append(messages, msgs...)
	}

	return messages
}

// Try to get a singl message from each pod specified
func GetSingleMessageFromEachPod(servers []string, path string) model.Message {
	chanMessages := make(chan model.Message, len(servers))

	var wg sync.WaitGroup
	for _, pod := range servers {
		wg.Add(1)
		go func(chanMessages chan<- model.Message, pod string) {
			defer wg.Done()
			adress := "https://" + pod + path

			fmt.Println("rq to", adress)
			//#nosec TODO avoid SSRF
			resp, err := http.Get(adress)
			if err != nil {
				fmt.Println("err", err)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("err", err)
				return
			}
			if !json.Valid(body) {
				return
			}

			msg := model.Message{}
			err = json.Unmarshal(body, &msg)
			if err != nil {
				fmt.Println("err", err)
				return
			}
			chanMessages <- msg

		}(chanMessages, pod)
	}

	var messages []model.Message

	wg.Wait()
	close(chanMessages)
	for msg := range chanMessages {
		messages = append(messages, msg)
	}

	if len(messages) == 0 {
		return model.Message{}
	}
	return messages[0]
}
