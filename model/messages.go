package model

import (
	"time"

	"github.com/ewenquim/horkruxes/service"
	"github.com/google/uuid"
)

type Message struct {
	// Stored and generated
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time

	// Stored and given by user
	DisplayedName   string // Name Chosen by author, no restriction but < 50 char
	AuthorBase64    string // Encoded in URL-safe Base64
	Content         string
	SignatureBase64 string // Encoded in URL-safe Base64
	MessageID       string // Used if the message is a comment to a publication

	// Only for display on client, computed from known values
	Correct       bool   `json:"-" gorm:"-"`
	Color         string `json:"-" gorm:"-"`
	DisplayedDate string `json:"-" gorm:"-"`
	Pod           string `gorm:"-"` // Not saved in db but tell where it is sent from so remains in JSON
}

func GetCommentsTo(s service.Service, messageID string) []Message {
	var messages []Message
	s.DB.Where("message_id = ?", messageID).Order("created_at desc").Find(&messages)
	return CleanMessagesOutFromDB(messages, s.GeneralConfig.URL)
}

// GetAllFromDB get data from db and checks some things
func GetAllFromDB(s service.Service) []Message {
	var messages []Message
	// s.DB.Where("correct = ?", true).Find(&messages)
	s.DB.Order("created_at desc").Find(&messages)
	return CleanMessagesOutFromDB(messages, s.GeneralConfig.URL)
}

// GetMessagesFromDB get data from db and checks some things
func GetMessagesFromDB(s service.Service) []Message {
	var messages []Message
	// s.DB.Where("correct = ?", true).Find(&messages)
	s.DB.Where("message_id IS NULL OR message_id=''").Order("created_at desc").Find(&messages)
	return CleanMessagesOutFromDB(messages, s.GeneralConfig.URL)
}

func GetMessagesFromAuthor(s service.Service, pubKeyBase64 string) []Message {
	var messages []Message
	s.DB.Where("author_base64 = ?", pubKeyBase64).Order("created_at desc").Find(&messages)
	return CleanMessagesOutFromDB(messages, s.GeneralConfig.URL)
}

func GetMostRecentMessage(s service.Service) Message {
	var message Message
	s.DB.Where("message_id IS NULL OR message_id=''").Order("created_at desc").First(&message)
	return CleanSingleMessageOutFromDB(message, s.GeneralConfig.URL)
}

func GetMostRecentComment(s service.Service, messageID string) Message {
	var message Message
	s.DB.Where("message_id IS NULL OR message_id=''").Where("message_id = ?", messageID).Order("created_at desc").First(&message)
	return CleanSingleMessageOutFromDB(message, s.GeneralConfig.URL)
}

func GetMessageFromDB(s service.Service, id string) (Message, error) {
	var message Message
	err := s.DB.First(&message, "id = ?", id).Error
	return CleanSingleMessageOutFromDB(message, s.GeneralConfig.URL), err
}

func NewMessage(s service.Service, message Message) error {
	if _, err := message.VerifyConditions(s); err != nil {
		return err
	}
	message.ID = uuid.NewString()
	message.CreatedAt = time.Now()
	return s.DB.Create(&message).Error
}
