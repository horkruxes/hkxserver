package exceptions

import "errors"

var ErrorRecordTooLongFound = errors.New("some fields of the message are too long")
