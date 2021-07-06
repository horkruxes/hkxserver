package service

import (
	"encoding/binary"
	"fmt"
	"html/template"
	"strings"

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

func ColorFromBytes(b []byte) string {
	if len(b) < 11 {
		return "red"
	}
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
