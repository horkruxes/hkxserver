package query

import "github.com/horkruxes/hkxserver/model"

// cleanMessagesOutFromDB get data from DB and do some checks and verifications
func cleanMessagesOutFromDB(messages []model.Message, url ...string) []model.Message {
	for i, message := range messages {
		messages[i] = cleanSingleMessageOutFromDB(message, url...)
	}
	return messages
}

func cleanSingleMessageOutFromDB(message model.Message, url ...string) model.Message {
	// url is set on server side and not re-set on client side
	if len(url) > 0 {
		message.Pod = url[0]
	}
	return message
}
