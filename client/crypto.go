package client

import (
	"crypto/ed25519"
	"encoding/base64"
)

type KeyGen struct {
	Pub            string
	NewPub         string
	ColorPrimary   string // Preview of the primary color of generated key pair
	ColorSecondary string // Preview of the secondary color of generated key pair
	Sec            string
	NewSec         string
	Sig            string
	DisplayedName  string
	Content        string
	MessageID      string
	Verif          bool
	Valid          bool
}

// GenKeys generates (cryptographically secured) a new pair of ed25519 keys
func GenKeys() KeyGen {
	pub, sec, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic("can't generate keys: the environment is not secure. On Linux, verify that /dev/urandom is available: " + err.Error())
	}
	pubString := base64.URLEncoding.EncodeToString(pub)
	secString := base64.URLEncoding.EncodeToString(sec)
	primary, secondary := colorsFromBase64(pubString)
	return KeyGen{NewPub: pubString, NewSec: secString, ColorPrimary: primary, ColorSecondary: secondary}
}
