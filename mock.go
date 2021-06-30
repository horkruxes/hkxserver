package main

import (
	"crypto/ed25519"

	"github.com/ewenquim/horkruxes-client/model"
)

func mock() model.PageData {
	// One time
	// Generates public and secret key
	pub, sec, _ := ed25519.GenerateKey(nil)
	pub2, sec2, _ := ed25519.GenerateKey(nil)
	// Sign every message TODO: possibility to do it in the site + advices
	// Needs the message, the pub and secret key. Outputs the signature
	// When the pod receives the message (content, public and signature) possibility to verify
	signature := ed25519.Sign(sec, []byte("hey guys, hello world"+string(pub)))
	signature2 := ed25519.Sign(sec, []byte("my first secure tweet"+string(pub)))
	signature3 := ed25519.Sign(sec2, []byte("lorem <strong>ipsum</strong>i skip\n lines"+string(pub2)))
	println("--------")
	return model.PageData{
		PageTitle: "Horkruxes",
		Messages: []model.Message{
			{DisplayedName: "ewen", AuthorPubKey: pub, Content: "hey guys, hello world", Signature: signature},
			{DisplayedName: "chloe.sa", AuthorPubKey: pub, Content: "my first secure tweet", Signature: signature2},
			{DisplayedName: "seraph", AuthorPubKey: []byte("2eb1ek2ed9g"), Content: `lorem https://ewen.quimerch.com/ <strong>ipsum</strong>i skip\n lines`},
			{DisplayedName: "marius", AuthorPubKey: pub2, Content: "lorem <strong>ipsum</strong>i skip\n lines", Signature: signature3},
		},
	}

}
