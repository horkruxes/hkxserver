package model

import (
	"sort"
)

// CleanMessagesOutFromDB get data from DB and do some checks and verifications
func CleanMessagesOutFromDB(messages []Message, url ...string) []Message {
	for i, message := range messages {
		messages[i] = CleanSingleMessageOutFromDB(message, url...)
	}
	return messages
}

func CleanSingleMessageOutFromDB(message Message, url ...string) Message {
	// url is set on server side and not re-set on client side
	if len(url) > 0 {
		message.Pod = url[0]
	}
	return message
}

func SortByDate(messages []Message) []Message {
	// Sort slice by date
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.After(messages[j].CreatedAt)
	})
	return messages
}
