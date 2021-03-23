package hw02_unpack_string //nolint:golint,stylecheck,revive

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var b strings.Builder
	var prevRune rune
	var hasPrevRune bool
	for _, r := range s {
		if unicode.IsDigit(r) {
			if !hasPrevRune {
				return "", ErrInvalidString
			}
			count, err := strconv.Atoi(string(r))
			if err != nil {
				return "", err
			}
			b.WriteString(strings.Repeat(string(prevRune), count))
			hasPrevRune = false
			continue
		}
		if hasPrevRune {
			b.WriteRune(prevRune)
		}
		prevRune = r
		hasPrevRune = true
	}
	if hasPrevRune {
		b.WriteRune(prevRune)
	}
	return b.String(), nil
}