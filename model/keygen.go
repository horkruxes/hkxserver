package model

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

type KeyGen struct {
	Pub           string
	NewPub        string
	Sec           string
	NewSec        string
	Sig           string
	DisplayedName string
	Content       string
	Verif         bool
	Valid         bool
}

func VerifyFromString(pub, sig, displayedName, msg string) bool {
	if pub == "" || sig == "" || displayedName == "" || msg == "" {
		return false
	}

	message := Message{
		AuthorBase64:    pub,
		DisplayedName:   displayedName,
		Content:         msg,
		SignatureBase64: sig,
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
func SignMessage(secBase64, pubBase64, displayedName, message string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	sec, _ := base64.StdEncoding.DecodeString(secBase64)
	pub, _ := base64.StdEncoding.DecodeString(pubBase64)
	fmt.Println("\n\n\nSIGNING", message, pubBase64, displayedName)

	msgToSign := append([]byte(message), pub...)
	msgToSign = append(msgToSign, []byte(displayedName)...)
	fmt.Println("msg 2 sign", msgToSign)

	signature := ed25519.Sign(sec, msgToSign)
	sigString := base64.StdEncoding.EncodeToString(signature)
	return sigString

}
