package client

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	"github.com/horkruxes/hkxserver/model"
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

// GenKeys generates (cryptographically secured) a new pair of ed25519 keys
func GenKeys() KeyGen {
	pub, sec, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic("can't generate keys: the environment is not secure. On Linux, verify that /dev/urandom is available: " + err.Error())
	}
	pubString := base64.URLEncoding.EncodeToString(pub)
	secString := base64.URLEncoding.EncodeToString(sec)
	primary, secondary := colorsFromBase64(pubString)
	return KeyGen{NewPub: pubString, NewSec: secString, ColorPrimary: primary, ColorSecondary: secondary}
}

func VerifyFromString(pub, sig, displayedName, msg, msgId string) bool {
	if pub == "" || sig == "" || displayedName == "" || msg == "" {
		return false
	}

	message := model.Message{
		AuthorBase64:    pub,
		DisplayedName:   displayedName,
		Content:         msg,
		SignatureBase64: sig,
		MessageID:       msgId,
	}

	return message.VerifyOwnerShip()
}

// SignMessage signs messages from base64 and return a base64 signature (empty string if the signature can't be generated)
// The signature contains these elements concatenated:
// The message (UTF-8 to bytes)
// The author's public key
// The author's declared name
// The eventual messageId, empty if this is an original post
func SignStrings(secBase64, pubBase64, displayedName, message, messageId string) string {
	sec, err := base64.URLEncoding.DecodeString(secBase64)
	if err != nil || len(sec) != ed25519.PrivateKeySize {
		return ""
	}
	pub, err := base64.URLEncoding.DecodeString(pubBase64)
	if err != nil {
		return ""
	}

	fmt.Println("\n\n\nSIGNING", message, pubBase64, displayedName, messageId)

	msgToSign := append([]byte(message), pub...)
	msgToSign = append(msgToSign, []byte(displayedName)...)
	msgToSign = append(msgToSign, []byte(messageId)...)
	fmt.Println("msg 2 sign", msgToSign)

	signature := ed25519.Sign(sec, msgToSign)
	sigString := base64.URLEncoding.EncodeToString(signature)
	return sigString
}

func SignMessage(msg model.Message, secretKey string) string {
	return SignStrings(secretKey, msg.AuthorBase64, msg.DisplayedName, msg.Content, msg.MessageID)
}
