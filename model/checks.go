package model

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

// VerifyConditions returns true if everything is ok
func (message Message) VerifyConditions() bool {
	return len(message.Content) < 5000 && len(message.DisplayedName) < 50
}

func (message Message) VerifyOwnerShip() bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	if message.AuthorPubKey == nil || message.Signature == nil || len(message.AuthorPubKey) == 0 || len(message.Signature) == 0 {
		return false
	}
	return ed25519.Verify(message.AuthorPubKey, []byte(message.Content+string(message.AuthorPubKey)), message.Signature)
}

// CleanMessagesOutFromDB get data from DB and do some checks and verifications
func CleanMessagesOutFromDB(messages []Message, url ...string) []Message {
	for i, message := range messages {
		messages[i] = CleanSingleMessageOutFromDB(message, url...)
	}
	return messages
}

func CleanSingleMessageOutFromDB(message Message, url ...string) Message {
	message.Correct = message.VerifyOwnerShip()
	message.AuthorBase64 = base64.StdEncoding.EncodeToString(message.AuthorPubKey)
	message.SignatureBase64 = base64.StdEncoding.EncodeToString(message.Signature)
	// url is set on server side and not re-set on client side
	if len(url) > 0 {
		message.Pod = url[0]
	}
	return message
}
