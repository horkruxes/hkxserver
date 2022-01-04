package model

import (
	"sort"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// Normalizes, escapes and verifies constraints
func (message *Message) Normalize() error {
	message.DisplayedName = strings.TrimSpace(message.DisplayedName)
	message.AuthorBase64 = strings.TrimSpace(message.AuthorBase64)
	message.Content = strings.TrimSpace(message.Content)
	message.SignatureBase64 = strings.TrimSpace(message.SignatureBase64)
	message.MessageID = strings.TrimSpace(message.MessageID)

	message.sanitizeAndEscapeHTML()

	return message.verifyConstraints()
}

// SanitizeAndEscapeHTML removes harmful HTML tags and escapes HTML characters
func (message *Message) sanitizeAndEscapeHTML() {
	policy := bluemonday.StrictPolicy()
	message.Content = policy.Sanitize(message.Content)
	message.DisplayedName = policy.Sanitize(message.DisplayedName)
}

func SortByDate(messages []Message) []Message {
	// Sort slice by date
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.After(messages[j].CreatedAt)
	})
	return messages
}
