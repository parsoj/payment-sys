package base62

import (
	"strings"
)

func Base62Encode(number int64) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	if number == 0 {
		return "0"
	}

	var result strings.Builder
	base := int64(len(charset))

	for number > 0 {
		remainder := number % base
		result.WriteByte(charset[remainder])
		number /= base
	}

	// The encoded result is in reverse order
	encoded := result.String()
	reverse := func(s string) string {
		runes := []rune(s)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	}

	return reverse(encoded)
}
