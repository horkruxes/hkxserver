package model

import (
	"time"

	"github.com/ewenquim/horkruxes/exceptions"
	"github.com/ewenquim/horkruxes/service"
	"github.com/google/uuid"
)

type Message struct {
	ID              string `gorm:"primary_key"`
	CreatedAt       time.Time
	DisplayedName   string `json:"displayedName"` // Name Chosen by author, no restriction but < 50 char
	Content         string `json:"content"`
	AuthorPubKey    []byte `json:"authorPubKey"`
	AuthorBase64    string `json:"authorBase64"`
	Signature       []byte `json:"signature"`
	SignatureBase64 string `json:"signatureBase64"`
	Correct         bool
	Color           string
	MessageID       string // Used if the message is a comment to a publication
	// Only for display, computed from known values
	AuthorURLSafe string `json:"authorURLSafe" gorm:"-"`
	DisplayedDate string `gorm:"-"`
	Pod           string `gorm:"-"`
}

func GetCommentsTo(s service.Service, messageID string) []Message {
	var messages []Message
	s.DB.Where("message_id = ?", messageID).Order("created_at desc").Find(&messages)
	return CleanMessagesOutFromDB(messages, s.ServerConfig.URL)
}

// GetMessagesFromDB get data from db and checks some things
func GetMessagesFromDB(s service.Service) []Message {
	var messages []Message
	// s.DB.Where("correct = ?", true).Find(&messages)
	s.DB.Where("message_id IS NULL OR message_id=''").Order("created_at desc").Find(&messages)
	return CleanMessagesOutFromDB(messages, s.ServerConfig.URL)
}

func GetMessagesFromAuthor(s service.Service, pubKeyBase64 string) []Message {
	var messages []Message
	// s.DB.Where("correct = ?", true).Find(&messages)
	s.DB.Where("author_base64 = ?", pubKeyBase64).Order("created_at desc").Find(&messages)

	return CleanMessagesOutFromDB(messages, s.ServerConfig.URL)
}

func GetMessageFromDB(s service.Service, id string) (Message, error) {
	var message Message
	err := s.DB.First(&message, "id = ?", id).Error
	return CleanSingleMessageOutFromDB(message, s.ServerConfig.URL), err
}

func NewMessage(s service.Service, message Message) error {
	if !message.VerifyConditions() {
		return exceptions.ErrorRecordTooLongFound
	}
	if !message.VerifyOwnerShip() {
		return exceptions.ErrorWrongSignature
	}
	message.Correct = true
	message.ID = uuid.NewString()
	message.CreatedAt = time.Now()
	return s.DB.Create(&message).Error
}
