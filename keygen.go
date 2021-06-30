package main

import (
	"crypto/ed25519"
	"fmt"
)

type KeyGen struct {
	Pub ed25519.PublicKey
	Sec ed25519.PrivateKey
}

func (message Message) verifyOwnerShip() bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	return ed25519.Verify(message.AuthorPubKey, []byte(message.Content+string(message.AuthorPubKey)), message.Signature)
}

func genKeys() KeyGen {
	pub, sec, _ := ed25519.GenerateKey(nil)
	return KeyGen{pub, sec}
}
