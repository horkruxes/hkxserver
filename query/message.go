package query

import (
	"time"

	"github.com/google/uuid"
	"github.com/horkruxes/hkxserver/model"
	"github.com/horkruxes/hkxserver/service"
)

func GetCommentsTo(s service.Service, messageID string) []model.Message {
	var messages []model.Message
	s.DB.Where("message_id = ?", messageID).Order("created_at desc").Find(&messages)
	return model.CleanMessagesOutFromDB(messages, s.GeneralConfig.URL)
}

// GetAll gets data from db and checks some things
func GetAll(s service.Service) []model.Message {
	var messages []model.Message
	// s.DB.Where("correct = ?", true).Find(&messages)
	s.DB.Order("created_at desc").Find(&messages)
	return model.CleanMessagesOutFromDB(messages, s.GeneralConfig.URL)
}

// GetMessages get data from db and checks some things
func GetMessages(s service.Service) []model.Message {
	var messages []model.Message
	// s.DB.Where("correct = ?", true).Find(&messages)
	s.DB.Where("message_id IS NULL OR message_id=''").Order("created_at desc").Find(&messages)
	return model.CleanMessagesOutFromDB(messages, s.GeneralConfig.URL)
}

func GetMessagesFromAuthor(s service.Service, pubKeyBase64 string) []model.Message {
	var messages []model.Message
	s.DB.Where("author_base64 = ?", pubKeyBase64).Order("created_at desc").Find(&messages)
	return model.CleanMessagesOutFromDB(messages, s.GeneralConfig.URL)
}

func GetMostRecentMessage(s service.Service) model.Message {
	var message model.Message
	s.DB.Where("message_id IS NULL OR message_id=''").Order("created_at desc").First(&message)
	return model.CleanSingleMessageOutFromDB(message, s.GeneralConfig.URL)
}

func GetMostRecentComment(s service.Service, messageID string) model.Message {
	var message model.Message
	s.DB.Where("message_id IS NULL OR message_id=''").Where("message_id = ?", messageID).Order("created_at desc").First(&message)
	return model.CleanSingleMessageOutFromDB(message, s.GeneralConfig.URL)
}

func GetMessage(s service.Service, id string) (model.Message, error) {
	var message model.Message
	err := s.DB.First(&message, "id = ?", id).Error
	return model.CleanSingleMessageOutFromDB(message, s.GeneralConfig.URL), err
}

// NewMessage sanitizes, checks validity and eligibilty and saves the message to the database
func NewMessage(s service.Service, message model.Message) (model.Message, error) {
	message.ID = uuid.NewString()
	message.CreatedAt = time.Now()
	if err := message.Normalize(true); err != nil {
		return model.Message{}, err
	}
	message.EscapesHTML()
	if err := VerifyServerConstraints(s, message); err != nil {
		return model.Message{}, err
	}
	err := s.DB.Create(&message).Error
	return message, err
}
