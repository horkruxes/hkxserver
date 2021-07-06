package model

import (
	"encoding/base64"
	"sort"

	"github.com/ewenquim/horkruxes/service"
)

// CleanMessagesOutFromDB get data from DB and do some checks and verifications
func CleanMessagesOutFromDB(messages []Message, url ...string) []Message {
	for i, message := range messages {
		messages[i] = CleanSingleMessageOutFromDB(message, url...)
	}
	return messages
}

func CleanSingleMessageOutFromDB(message Message, url ...string) Message {
	message.Correct = message.VerifyOwnerShip()
	message.AuthorBase64 = base64.StdEncoding.EncodeToString(message.AuthorPubKey)
	message.SignatureBase64 = base64.StdEncoding.EncodeToString(message.Signature)
	// url is set on server side and not re-set on client side
	if len(url) > 0 {
		message.Pod = url[0]
	}
	return message
}

// CleanMessagesOutFromDB get data from DB and do some checks and verifications
func CleanMessagesClientSide(messages []Message) []Message {
	for i, message := range messages {
		message = CleanSingleMessageOutFromDB(message)
		messages[i] = CleanSingleMessageClientSide(message)
	}
	return messages
}

func CleanSingleMessageClientSide(message Message) Message {
	message.AuthorURLSafe = service.Base64ToSafeURL(message.AuthorBase64)
	message.DisplayedDate = message.CreatedAt.Format("2 Jan 2006 15:04")
	message.Color = service.ColorFromBytes(message.AuthorPubKey)
	return message
}

func SortByDate(messages []Message) []Message {
	// Sort slice by date
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.After(messages[j].CreatedAt)
	})
	return messages
}
