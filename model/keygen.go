package model

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

type KeyGen struct {
	Pub     string
	NewPub  string
	Sec     string
	NewSec  string
	Sig     string
	Content string
	Verif   bool
	Valid   bool
}

func VerifyFromString(pub, sig, msg string) bool {
	if pub == "" || sig == "" || msg == "" {
		return false
	}
	pubByte, err := base64.StdEncoding.DecodeString(pub)
	if err != nil {
		return false
	}
	sigByte, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return false
	}
	message := Message{
		AuthorPubKey: pubByte,
		Content:      msg,
		Signature:    sigByte,
	}

	return message.VerifyOwnerShip()
}

func GenKeys() KeyGen {
	pub, sec, _ := ed25519.GenerateKey(nil)
	pubString := base64.StdEncoding.EncodeToString(pub)
	secString := base64.StdEncoding.EncodeToString(sec)
	return KeyGen{NewPub: pubString, NewSec: secString}
}

// SignMessage signs messages from base64 and return a base64 signature
func SignMessage(secString, pubString, message string) string {
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
