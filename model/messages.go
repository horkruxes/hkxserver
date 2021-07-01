package model

import (
	"github.com/ewenquim/horkruxes/exceptions"
	"github.com/ewenquim/horkruxes/service"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	DisplayedName   string `json:"displayedName"` // Name Chosen by author, no restriction but < 50 char
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
	if !message.VerifyConditions() {
		return exceptions.ErrorRecordTooLongFound
	}
	if message.VerifyOwnerShip() {
		message.Correct = true
		return s.DB.Create(&message).Error
	}
	return nil
}
