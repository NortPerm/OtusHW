package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func addRuneToBuilder(b *strings.Builder, r rune, times int32) error {
	for times > 0 {
		b.WriteRune(r)
		times = times - 1
	}
	return nil
}

func escapeRune(r rune) bool {
	return r == '\\' || unicode.IsDigit(r)

}

func Unpack(input string) (string, error) {
	lastRune := rune(0)
	escapeMode := false
	runes := []rune(input + "\x00") // add terminator zero rune
	var result strings.Builder
	for _, currentRune := range runes {
		if currentRune == '\\' && !escapeMode {
			escapeMode = true
		} else {
			if !escapeRune(currentRune) && escapeMode {
				return "", ErrInvalidString
			}
			if unicode.IsDigit(currentRune) && !escapeMode {
				if lastRune > 0 {
					addRuneToBuilder(&result, lastRune, currentRune-'0')
				} else {
					return "", ErrInvalidString
				}
				lastRune = 0
			} else {
				if lastRune > 0 {
					result.WriteRune(lastRune)
				}
				lastRune = currentRune
			}
			escapeMode = false
		}

	}
	return result.String(), nil
}
