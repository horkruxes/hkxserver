package model

import (
	"crypto/ed25519"
	"fmt"
)

// VerifyConditions returns true if everything is ok
func (message Message) VerifyConditions() bool {
	return len(message.Content) > 140 && len(message.Content) < 5000 && len(message.DisplayedName) < 50
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
