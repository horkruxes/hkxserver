package model

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

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

// SignMessage signs messages from base64 and return a base64 signature (empty string if the signature can't be generated)
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
