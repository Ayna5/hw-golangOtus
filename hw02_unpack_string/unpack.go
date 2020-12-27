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
		switch r {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			count, err := strconv.Atoi(string(r))
			if err != nil {
				return "", err
			}
			if unicode.IsDigit(prevRune) {
				return "", ErrInvalidString
			}
			if count > 0 {
				b.WriteString(strings.Repeat(string(prevRune), count))
			}
		default:
			if !unicode.IsDigit(prevRune) && prevRune != 0 {
				b.WriteRune(prevRune)
			}
			if strings.HasSuffix(s, string(r)) && i == len(s)-len(string(r)) {
				b.WriteRune(r)
			}
		}
		prevRune = r
	}
	return b.String(), nil
}
