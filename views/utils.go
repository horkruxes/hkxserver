package views

import (
	"encoding/binary"
	"fmt"
	"html/template"
	"sort"
	"strings"

	"github.com/ewenquim/horkruxes/model"
	"github.com/russross/blackfriday"
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
	// Sort slice by date
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.After(messages[j].CreatedAt)
	})
	return messages
}

func ColorFromBytes(b []byte) string {
	red := int(binary.BigEndian.Uint32(b[:]))%223 + 16
	green := int(binary.BigEndian.Uint32(b[10:]))%223 + 16
	blue := int(binary.BigEndian.Uint32(b[5:]))%223 + 16
	return fmt.Sprintf("#%X%X%X", red, green, blue)
	// hue := int(binary.BigEndian.Uint32(b)) % 360
	// light := ((hue+int(binary.BigEndian.Uint32(b[10:])))%3 + 1) * 25
	// fmt.Printf("hsl(%v, 100%%, %v%%)", hue, light)
	// return fmt.Sprintf("hsl(%v, 100%%, %v%%)", hue, light)
}

func MarkDowner(args ...interface{}) template.HTML {
	s := blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}
