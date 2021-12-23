package client

import (
	"encoding/base64"
	"fmt"
	"hash/fnv"

	"github.com/horkruxes/hkxserver/model"
)

// getRandomNumber returns a seemingly random but deterministic number between 16 & 239=255-16.
// It does not uses a cryptographically secure hash for the moment
func getRandomNumber(b []byte) uint8 {
	h := fnv.New32a()
	//#nosec
	h.Write(b)
	return uint8(h.Sum32()%223 + 16)
}

func colorsFromBase64(name string) (string, string) {
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

// CleanMessagesOutFromDB get data from DB and do some checks and verifications
func CleanMessagesClientSide(messages []model.Message) []model.Message {
	for i, message := range messages {
		message = model.CleanSingleMessageOutFromDB(message)
		messages[i] = CleanSingleMessageClientSide(message)
	}
	return messages
}

func CleanSingleMessageClientSide(message model.Message) model.Message {
	message.DisplayedDate = message.CreatedAt.Format("2 Jan 2006 15:04")

	message.ColorPrimary, message.ColorSecondary = colorsFromBase64(message.AuthorBase64)
	message.Correct = message.VerifyOwnerShip()
	return message
}
