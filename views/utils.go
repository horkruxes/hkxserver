package views

import (
	"strings"

	"github.com/ewenquim/horkruxes-client/model"
)

func Base64ToSafeURL(s string) string {
	s = strings.ReplaceAll(s, "+", ".")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, "=", "_")
	return s
}

func SafeURLToBase64(s string) string {
	s = strings.ReplaceAll(s, ".", "+")
	s = strings.ReplaceAll(s, "-", "/")
	s = strings.ReplaceAll(s, "_", "=")
	return s
}

// CleanMessagesOutFromDB get data from DB and do some checks and verifications
func CleanMessagesClientSide(messages []model.Message) []model.Message {
	messages = model.CleanMessagesOutFromDB(messages)
	for i, message := range messages {
		messages[i].AuthorURLSafe = Base64ToSafeURL(message.AuthorBase64)
		messages[i].DisplayedDate = message.CreatedAt.Format("2 Jan 2006 15:04")
	}
	return messages
}
