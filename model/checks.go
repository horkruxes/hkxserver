package model

import (
	"crypto/ed25519"
	"encoding/base64"

	"github.com/horkruxes/hkxserver/exceptions"
)

// VerifyConstraints returns HTTP status code and an error
// Checks that the messages constraints are inherently met
// -independently from the database & server.
func (message Message) VerifyConstraints() error {
	if len(message.Content) > 50000 || len(message.DisplayedName) > 50 {
		return exceptions.ErrorFieldsTooLong
	} else if len(message.Content) < 140 {
		return exceptions.ErrorContentTooShort
	} else if !message.VerifyOwnerShip() {
		return exceptions.ErrorWrongSignature
	}
	return nil
}

func (message Message) VerifyOwnerShip() bool {
	pubBytes, err := base64.URLEncoding.DecodeString(message.AuthorBase64)
	if err != nil || len(pubBytes) != ed25519.PublicKeySize {
		return false
	}
	sigBytes, err := base64.URLEncoding.DecodeString(message.SignatureBase64)
	if err != nil {
		return false
	}
	if len(pubBytes) == 0 || len(sigBytes) == 0 {
		return false
	}
	messageWithInfo := append([]byte(message.Content), pubBytes...)
	messageWithInfo = append(messageWithInfo, []byte(message.DisplayedName)...)
	messageWithInfo = append(messageWithInfo, []byte(message.MessageID)...)
	// fmt.Println("\n\n\nVERIFYING", message.Content[:20], message.AuthorBase64, message.DisplayedName, message.MessageID)
	// fmt.Println(messageWithInfo)

	return ed25519.Verify(pubBytes, messageWithInfo, sigBytes)
}
