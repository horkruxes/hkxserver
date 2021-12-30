package model

import (
	"time"
)

type Message struct {
	// Stored and generated
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time

	// Stored and given by user
	DisplayedName   string `form:"title" json:"title"`           // Name Chosen by author, no restriction but < 50 char
	AuthorBase64    string `form:"public-key" json:"public-key"` // Encoded in URL-safe Base64
	Content         string `form:"message" json:"message"`
	SignatureBase64 string `form:"signature" json:"signature"` // Encoded in URL-safe Base64
	MessageID       string `form:"answer-to" json:"answer-to"` // Used if the message is a comment to a publication

	// Only for display on client, computed from known values
	Correct        bool   `json:"-" gorm:"-"`
	ColorPrimary   string `json:"-" gorm:"-"`
	ColorSecondary string `json:"-" gorm:"-"`
	DisplayedDate  string `json:"-" gorm:"-"`
	Pod            string `gorm:"-"` // Not saved in db but tell where it is sent from so remains in JSON
}
