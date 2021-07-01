package model

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	"github.com/ewenquim/horkruxes/service"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	DisplayedName   string `json:"displayedName"` // Name Chosen by author, bo restriction
	Content         string `json:"content"`
	AuthorPubKey    []byte `json:"authorPubKey"`
	AuthorBase64    string `json:"authorBase64"`
	Signature       []byte `json:"signature"`
	SignatureBase64 string `json:"signatureBase64"`
	Correct         bool
	Color           string
	// Only for display, computed from known values
	AuthorURLSafe string `json:"authorURLSafe" gorm:"-"`
	DisplayedDate string `gorm:"-"`
	Pod           string `gorm:"-"`
}

// GetMessagesFromDB get data from db and checks some things
func GetMessagesFromDB(s service.Service) []Message {
	var messages []Message
	// s.DB.Where("correct = ?", true).Find(&messages)
	s.DB.Order("created_at desc").Find(&messages)
	return CleanMessagesOutFromDB(messages, s.ServerConfig.URL)
}

func GetMessagesFromAuthor(s service.Service, pubKeyBase64 string) []Message {
	var messages []Message
	// s.DB.Where("correct = ?", true).Find(&messages)
	s.DB.Where("author_base64 = ?", pubKeyBase64).Order("created_at desc").Find(&messages)

	return CleanMessagesOutFromDB(messages, s.ServerConfig.URL)
}

func GetMessageFromDB(s service.Service, id string) Message {
	var message Message
	s.DB.Find(&message, id)
	return message
}

func NewMessage(s service.Service, message *Message) error {
	if message.VerifyOwnerShip() {
		message.Correct = true
		return s.DB.Create(&message).Error
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
func CleanMessagesOutFromDB(messages []Message, url ...string) []Message {
	for i, message := range messages {
		messages[i].Correct = message.VerifyOwnerShip()
		messages[i].AuthorBase64 = base64.StdEncoding.EncodeToString(message.AuthorPubKey)
		messages[i].Color = ColorFromString(string(message.AuthorPubKey))
		messages[i].SignatureBase64 = base64.StdEncoding.EncodeToString(message.Signature)
		// url is set on server side and not re-set on client side
		if len(url) > 0 {
			messages[i].Pod = url[0]
		}
	}
	return messages
}
