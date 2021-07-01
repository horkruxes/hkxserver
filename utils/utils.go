package utils

import "strings"

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