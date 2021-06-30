package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

type KeyGen struct {
	Pub string
	Sec string
	Sig string
}

func (message Message) verifyOwnerShip() bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	return ed25519.Verify(message.AuthorPubKey, []byte(message.Content+string(message.AuthorPubKey)), message.Signature)
}

func genKeys() KeyGen {
	pub, sec, _ := ed25519.GenerateKey(nil)
	pubString := base64.StdEncoding.EncodeToString(pub)
	secString := base64.StdEncoding.EncodeToString(sec)
	return KeyGen{pubString, secString, ""}
}

// signMessage signs messages from base64 and return a base64 signature
func signMessage(secString, pubString, message string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	sec, _ := base64.StdEncoding.DecodeString(secString)
	pub, _ := base64.StdEncoding.DecodeString(pubString)
	signature := ed25519.Sign(sec, []byte(message+string(pub)))
	sigString := base64.StdEncoding.EncodeToString(signature)
	return sigString

}
