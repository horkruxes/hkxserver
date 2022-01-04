package model

import (
	"crypto/ed25519"
	"encoding/base64"
	"strings"

	"github.com/google/uuid"
	"github.com/horkruxes/hkxserver/exceptions"
	"github.com/microcosm-cc/bluemonday"
)

// Verify that the message is sanitized and has a valid signature
func (message Message) Verify() error {
	if err := message.verifyConstraints(); err != nil {
		return err
	} else if !message.isSanitized() {
		return exceptions.ErrorContentWithHTML
	} else if !message.verifyOwnerShip() {
		return exceptions.ErrorWrongSignature
	}
	return nil
}

func (message Message) isSanitized() bool {
	if !isEscaped(message.DisplayedName) || !isEscaped(message.Content) {
		return false
	}
	policy := bluemonday.StrictPolicy()
	if message.Content != policy.Sanitize(message.Content) {
		return false
	}
	if message.DisplayedName != policy.Sanitize(message.DisplayedName) {
		return false
	}
	return true
}

func isEscaped(s string) bool {
	return !strings.ContainsAny(s, "<>\"'")
}

// VerifyConstraints returns an error.
// Checks that the messages constraints are inherently met
// -independently from the database & server.
func (message Message) verifyConstraints() error {
	if len(message.Content) > 50000 || len(message.DisplayedName) > 50 {
		return exceptions.ErrorFieldsTooLong
	} else if len(message.Content) < 140 {
		return exceptions.ErrorContentTooShort
	} else if _, err := uuid.Parse(message.ID); message.ID != "" && err != nil {
		return err
	}
	return nil
}

// VerifyOwnerShip returns true if the signature is valid.
func (message Message) verifyOwnerShip() bool {
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
