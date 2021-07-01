package views

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/ewenquim/horkruxes/model"
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
		messages[i].Color = ColorFromBytes(message.AuthorPubKey)

	}
	return messages
}

func ColorFromBytes(b []byte) string {

	// java String#hashCode
	red := int(binary.BigEndian.Uint32(b[:])) % 255
	green := int(binary.BigEndian.Uint32(b[10:])) % 255
	blue := int(binary.BigEndian.Uint32(b[5:])) % 255
	fmt.Println("hi")

	return fmt.Sprintf("#%X%X%X", red, green, blue)
}
