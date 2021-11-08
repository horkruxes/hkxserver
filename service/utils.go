package service

import (
	"encoding/base64"
	"fmt"
	"hash/fnv"
	"html/template"

	"github.com/russross/blackfriday/v2"
)

// getRandomNumber returns a seemingly random but deterministic number between 16 & 239=255-16.
// It does not uses a cryptographically secure hash for the moment
func getRandomNumber(b []byte) uint8 {
	h := fnv.New32a()
	h.Write(b)
	return uint8(h.Sum32()%223 + 16)
}

func ColorsFromBase64(name string) (string, string) {
	b, err := base64.URLEncoding.DecodeString(name)
	if err != nil {
		return "red", "red"
	}

	// b must be 32 bytes long
	if len(b) < 32 {
		return "red", "red"
	}

	red := getRandomNumber(b[:5])
	green := getRandomNumber(b[5:10])
	blue := getRandomNumber(b[10:16])

	redSec := getRandomNumber(b[16:21])
	greenSec := getRandomNumber(b[21:26])
	blueSec := getRandomNumber(b[26:32])

	return fmt.Sprintf("#%X%X%X", red, green, blue), fmt.Sprintf("#%X%X%X", redSec, greenSec, blueSec)
	// hue := int(binary.BigEndian.Uint32(b)) % 360
	// light := ((hue+int(binary.BigEndian.Uint32(b[10:])))%3 + 1) * 25
	// fmt.Printf("hsl(%v, 100%%, %v%%)", hue, light)
	// return fmt.Sprintf("hsl(%v, 100%%, %v%%)", hue, light)
}

func MarkDowner(args ...interface{}) template.HTML {
	s := blackfriday.Run([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}
