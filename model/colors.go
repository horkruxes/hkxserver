package model

import (
	"math"
	"strconv"
	"strings"
)

func ColorFromString(s string) string {
	// java String#hashCode
	str := strings.ToLower(s)
	hash := 0
	for i := 0; i < len(str); i++ {
		num, _ := strconv.Atoi(string(str[i]))
		hash = num + ((hash << 5) - hash)
	}
	base := int(math.Abs(float64(hash)))

	colors := []string{
		"pink",
		"#9b88ee",
		"GainsBoRo",
		"yellowGreen",
		"skyBlue",
		"salmon",
		"sandyBrown",
		"paleGreen",
		"paleTurquoise",
		"red",
	}

	return "fill:" + strings.ToLower(colors[trueMod(base, len(colors))])
}

func trueMod(n int, m int) int {
	return ((n % m) + m) % m
}
