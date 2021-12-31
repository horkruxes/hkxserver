package model

import (
	"crypto/ed25519"
	"encoding/base64"
	"strings"

	"github.com/google/uuid"
	"github.com/horkruxes/hkxserver/exceptions"
	"github.com/microcosm-cc/bluemonday"
)

// Normalizes, verifies integrity and optionally verifies the message signature
func (message *Message) Normalize(verifySignature bool) error {
	message.DisplayedName = strings.TrimSpace(message.DisplayedName)
	message.AuthorBase64 = strings.TrimSpace(message.AuthorBase64)
	message.Content = strings.TrimSpace(message.Content)
	message.SignatureBase64 = strings.TrimSpace(message.SignatureBase64)
	message.MessageID = strings.TrimSpace(message.MessageID)

	return message.verifyConstraints(verifySignature)
}

func (message *Message) EscapesHTML() {
	message.Content = bluemonday.UGCPolicy().Sanitize(message.Content)
	message.DisplayedName = bluemonday.StrictPolicy().Sanitize(message.DisplayedName)
}

// VerifyConstraints returns an error.
// Checks that the messages constraints are inherently met
// -independently from the database & server.
func (message Message) verifyConstraints(verifySignature bool) error {
	if len(message.Content) > 50000 || len(message.DisplayedName) > 50 {
		return exceptions.ErrorFieldsTooLong
	} else if len(message.Content) < 140 {
		return exceptions.ErrorContentTooShort
	} else if _, err := uuid.Parse(message.ID); err != nil {
		return err
	} else if verifySignature && !message.verifyOwnerShip() {
		return exceptions.ErrorWrongSignature
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
