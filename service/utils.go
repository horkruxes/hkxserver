package service

import (
	"encoding/binary"
	"fmt"
	"html/template"

	"github.com/russross/blackfriday"
)

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
	s := blackfriday.MarkdownBasic([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}
