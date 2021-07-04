package views

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ewenquim/horkruxes/model"
)

type Resp struct {
	Response []model.Message `json:"response"`
}

func getMessagesFrom(path string, urls ...string) []model.Message {
	var messages []model.Message
	for _, url := range urls {
		adress := "https://" + url + path

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

		msg := &Resp{}
		err = json.Unmarshal(body, msg)
		if err != nil {
			fmt.Println("err", err)
			continue
		}
		messages = append(messages, (*msg).Response...)
	}
	return messages
}
