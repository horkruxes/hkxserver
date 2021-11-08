package model

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	"github.com/ewenquim/horkruxes/service"
)

type KeyGen struct {
	Pub            string
	NewPub         string
	ColorPrimary   string // Preview of the primary color of generated key pair
	ColorSecondary string // Preview of the secondary color of generated key pair
	Sec            string
	NewSec         string
	Sig            string
	DisplayedName  string
	Content        string
	MessageID      string
	Verif          bool
	Valid          bool
}

func VerifyFromString(pub, sig, displayedName, msg, msgId string) bool {
	if pub == "" || sig == "" || displayedName == "" || msg == "" {
		return false
	}

	message := Message{
		AuthorBase64:    pub,
		DisplayedName:   displayedName,
		Content:         msg,
		SignatureBase64: sig,
		MessageID:       msgId,
	}

	return message.VerifyOwnerShip()
}

// GenKeys generates (cryptographically secured) a new pair of ed25519 keys
func GenKeys() KeyGen {
	pub, sec, _ := ed25519.GenerateKey(nil)
	pubString := base64.URLEncoding.EncodeToString(pub)
	secString := base64.URLEncoding.EncodeToString(sec)
	primary, secondary := service.ColorsFromBase64(pubString)
	return KeyGen{NewPub: pubString, NewSec: secString, ColorPrimary: primary, ColorSecondary: secondary}
}

// SignMessage signs messages from base64 and return a base64 signature
// The signature contains these elements concatenated:
// The message (UTF-8 to bytes)
// The author's public key
// The author's declared name
// The eventual messageId, empty if this is an original post
func SignMessage(secBase64, pubBase64, displayedName, message, messageId string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	sec, _ := base64.URLEncoding.DecodeString(secBase64)
	pub, _ := base64.URLEncoding.DecodeString(pubBase64)

	fmt.Println("\n\n\nSIGNING", message, pubBase64, displayedName, messageId)

	msgToSign := append([]byte(message), pub...)
	msgToSign = append(msgToSign, []byte(displayedName)...)
	msgToSign = append(msgToSign, []byte(messageId)...)
	fmt.Println("msg 2 sign", msgToSign)

	signature := ed25519.Sign(sec, msgToSign)
	sigString := base64.URLEncoding.EncodeToString(signature)
	return sigString

}
