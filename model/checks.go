package model

import (
	"crypto/ed25519"
	"fmt"
	"time"

	"github.com/ewenquim/horkruxes/exceptions"
	"github.com/ewenquim/horkruxes/service"
	"github.com/gofiber/fiber/v2"
)

// VerifyConditions returns HTTP status code and an error
func (message Message) VerifyConditions(s service.Service) (int, error) {
	if len(message.Content) > 50000 || len(message.DisplayedName) > 50 {
		return fiber.StatusBadRequest, exceptions.ErrorFieldsTooLong
	} else if len(message.Content) < 140 {
		return fiber.StatusBadRequest, exceptions.ErrorContentTooShort
	} else if !message.VerifyOwnerShip() {
		return fiber.StatusBadRequest, exceptions.ErrorWrongSignature
	} else {
		lastPost := GetMostRecentMessage(s).CreatedAt
		if time.Since(lastPost) < time.Hour && !message.authorTrusted(s) {
			return fiber.StatusNotAcceptable, exceptions.ErrorTooSoonUnregistered
		} else if time.Since(lastPost) < 5*time.Minute && message.authorTrusted(s) {
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
	if message.AuthorPubKey == nil || message.Signature == nil || len(message.AuthorPubKey) == 0 || len(message.Signature) == 0 {
		return false
	}
	return ed25519.Verify(message.AuthorPubKey, []byte(message.Content+string(message.AuthorPubKey)), message.Signature)
}

func (message Message) authorTrusted(s service.Service) bool {
	for _, key := range s.ServerConfig.TrustedKeys {
		if key == message.AuthorBase64 {
			return true
		}
	}
	return false
}
