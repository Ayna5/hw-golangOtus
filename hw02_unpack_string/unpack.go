package hw02_unpack_string //nolint:golint,stylecheck

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
	if len(s) > 0 && unicode.IsDigit([]rune(s)[0]) {
		return "", ErrInvalidString
	}
	for i, r := range s {
		//nolint:nestif // it should be refactored when I'll do task with asterisk
		if unicode.IsDigit(r) {
			if unicode.IsDigit(prevRune) {
				return "", ErrInvalidString
			}
			count, _ := strconv.Atoi(string(r))
			if count > 0 {
				b.WriteString(strings.Repeat(string(prevRune), count))
			}
		} else {
			if !unicode.IsDigit(prevRune) && prevRune != 0 {
				b.WriteRune(prevRune)
			}
			if i == len(s)-len(string(r)) {
				b.WriteRune(r)
			}
		}
		prevRune = r
	}
	return b.String(), nil
}
