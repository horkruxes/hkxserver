package query

import (
	"time"

	"github.com/horkruxes/hkxserver/exceptions"
	"github.com/horkruxes/hkxserver/model"
	"github.com/horkruxes/hkxserver/service"
)

func VerifyServerConstraints(s service.Service, message model.Message) error {
	if html := s.ContentPolicy.Sanitize(message.Content); html != message.Content {
		return exceptions.ErrorContentWithHTML
	}
	if err := VerifyTiming(s, message); err != nil {
		return err
	}
	return nil
}

func VerifyTiming(s service.Service, message model.Message) error {
	var lastPost time.Time
	if message.MessageID == "" {
		lastPost = GetMostRecentMessage(s).CreatedAt
	} else {
		lastPost = GetMostRecentComment(s, message.MessageID).CreatedAt
	}
	trusted := authorTrusted(s, message)

	if !trusted && time.Since(lastPost) < time.Hour {
		return exceptions.ErrorTooSoonUnregistered
	} else if trusted && time.Since(lastPost) < 30*time.Second {
		return exceptions.ErrorTooSoonRegistered
	}
	return nil
}

func authorTrusted(s service.Service, message model.Message) bool {
	for _, key := range s.ServerConfig.TrustedKeys {
		if key == message.AuthorBase64 {
			return true
		}
	}
	return false
}
