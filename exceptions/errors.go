package exceptions

import "errors"

var ErrorFieldsTooLong = errors.New("some fields of the message are too long. Name must be max 50 characters and content 50k characters")
var ErrorContentTooShort = errors.New("the content of your message is less than 140 characters long. Please write meaningful and informative content")
var ErrorWrongSignature = errors.New("wrong signature")
var ErrorTooSoonUnregistered = errors.New("when not registered on a pod, you must wait 1 hour after the last message posted on the server to avoid spam")
var ErrorTooSoonRegistered = errors.New("when registered on a pod, you must wait 30 sec after the last message posted on the server to avoid spam")
