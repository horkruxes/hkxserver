package exceptions

import "errors"

var ErrorRecordTooLongFound = errors.New("some fields of the message are too long")
var WrongSignature = errors.New("wrong signature")
