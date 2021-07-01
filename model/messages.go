package model

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	DisplayedName   string `json:"displayedName"` // Name Chosen by author, bo restriction
	Content         string `json:"content"`
	AuthorPubKey    []byte `json:"authorPubKey"`
	AuthorBase64    string `json:"authorBase64"`
	AuthorURLSafe   string `json:"authorURLSafe"`
	Signature       []byte `json:"signature"`
	SignatureBase64 string `json:"signatureBase64"`
	Correct         bool
	Color           string
	// Date          time.Time
}

// GetMessagesFromDB get data from db and checks some things
func GetMessagesFromDB(db *gorm.DB) []Message {
	var messages []Message
	// db.Where("correct = ?", true).Find(&messages)
	db.Order("created_at desc").Find(&messages)
	return CleanMessagesOutFromDB(messages)
}

func GetMessagesFromAuthor(db *gorm.DB, pubKeyBase64 string) []Message {
	var messages []Message
	// db.Where("correct = ?", true).Find(&messages)
	db.Where("author_base64 = ?", pubKeyBase64).Order("created_at desc").Find(&messages)

	return CleanMessagesOutFromDB(messages)
}

func GetMessageFromDB(db *gorm.DB, id string) Message {
	var message Message
	db.Find(&message, id)
	return message
}

func NewMessage(db *gorm.DB, message *Message) error {
	if message.VerifyOwnerShip() {
		message.Correct = true
		return db.Create(&message).Error
	}
	return nil
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
func CleanMessagesOutFromDB(messages []Message) []Message {
	for i, message := range messages {
		messages[i].Correct = message.VerifyOwnerShip()
		messages[i].AuthorBase64 = base64.StdEncoding.EncodeToString(message.AuthorPubKey)
		messages[i].Color = ColorFromString(string(message.AuthorPubKey))
		messages[i].SignatureBase64 = base64.StdEncoding.EncodeToString(message.Signature)
	}
	return messages
}
