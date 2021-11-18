package model

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/horkruxes/hkxserver/exceptions"
	"github.com/horkruxes/hkxserver/service"
)

// VerifyConditions returns HTTP status code and an error
func (message Message) VerifyConditions(s service.Service) (int, error) {
	if s.ServerConfig.Debug {
		return fiber.StatusAccepted, nil
	}

	if len(message.Content) > 50000 || len(message.DisplayedName) > 50 {
		return fiber.StatusBadRequest, exceptions.ErrorFieldsTooLong
	} else if len(message.Content) < 140 {
		return fiber.StatusBadRequest, exceptions.ErrorContentTooShort
	} else if html := s.ContentPolicy.Sanitize(message.Content); html == message.Content {
		return fiber.StatusBadRequest, exceptions.ErrorContentWithHTML
	} else if !message.VerifyOwnerShip() {
		return fiber.StatusBadRequest, exceptions.ErrorWrongSignature
	} else {
		var lastPost time.Time
		if message.MessageID == "" {
			lastPost = GetMostRecentMessage(s).CreatedAt
		} else {
			lastPost = GetMostRecentComment(s, message.MessageID).CreatedAt
		}
		trusted := message.authorTrusted(s)
		if !trusted && time.Since(lastPost) < time.Hour {
			return fiber.StatusNotAcceptable, exceptions.ErrorTooSoonUnregistered
		} else if trusted && time.Since(lastPost) < 30*time.Second {
			return fiber.StatusNotAcceptable, exceptions.ErrorTooSoonRegistered
		}
	}
	return fiber.StatusAccepted, nil
}

func (message Message) VerifyOwnerShip() bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	pubBytes, err := base64.URLEncoding.DecodeString(message.AuthorBase64)
	if err != nil {
		return false
	}
	sigBytes, err := base64.URLEncoding.DecodeString(message.SignatureBase64)
	if err != nil {
		return false
	}
	if len(pubBytes) == 0 || len(sigBytes) == 0 {
		return false
	}
	fmt.Println("\n\n\nVERIFYING", message.Content[:20], message.AuthorBase64, message.DisplayedName, message.MessageID)
	messageWithInfo := append([]byte(message.Content), pubBytes...)
	messageWithInfo = append(messageWithInfo, []byte(message.DisplayedName)...)
	messageWithInfo = append(messageWithInfo, []byte(message.MessageID)...)
	fmt.Println("msg 2 verify", messageWithInfo)

	return ed25519.Verify(pubBytes, messageWithInfo, sigBytes)
}

func (message Message) authorTrusted(s service.Service) bool {
	for _, key := range s.ServerConfig.TrustedKeys {
		if key == message.AuthorBase64 {
			return true
		}
	}
	return false
}
