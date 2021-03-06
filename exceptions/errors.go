package exceptions

import "errors"

var (
	ErrorFieldsTooLong       = errors.New("some fields of the message are too long. Name must be max 50 characters and content 50k characters")
	ErrorContentTooShort     = errors.New("the content of your message is less than 140 characters long. Please write meaningful and informative content")
	ErrorContentWithHTML     = errors.New("the content of your message contains HTML. Please write plain text")
	ErrorWrongSignature      = errors.New("wrong signature")
	ErrorTooSoonUnregistered = errors.New("when not registered on a pod, you must wait 1 hour after the last message posted on the server to avoid spam")
	ErrorTooSoonRegistered   = errors.New("when registered on a pod, you must wait 30 sec after the last message posted on the server to avoid spam")
)
