package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/horkruxes/hkxserver/model"
	"github.com/horkruxes/hkxserver/service"
)

func GetMessagesFrom(s service.Service, path string) []model.Message {
	chanMessages := make(chan []model.Message, len(path))
	chanFinal := make(chan []model.Message, len(path))

	var wg sync.WaitGroup
	for _, pod := range s.GeneralConfig.TrustedPods {
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

	go func(chanMessages <-chan []model.Message, chanFinal chan<- []model.Message) {
		for msgs := range chanMessages {
			messages = append(messages, msgs...)
		}

	}(chanMessages, chanFinal)
	wg.Wait()
	close(chanMessages)

	return messages
}

func GetSingleMessageFromEachPod(s service.Service, path string) []model.Message {
	var messages []model.Message
	for _, pod := range s.GeneralConfig.TrustedPods {
		adress := "https://" + pod + path

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
	return messages
}
