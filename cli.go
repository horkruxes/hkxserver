package main

import (
	"fmt"
	"net/http"

	"github.com/inconshreveable/go-update"
)

func doUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Can't fetch latest upgrade online")
		return err
	}
	defer resp.Body.Close()
	fmt.Println("Applying update")
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		// error handling
		fmt.Println("Can't update")
		return err
	}
	return err
}

func version() {
	fmt.Println("Horkruxes v0.1")
}
