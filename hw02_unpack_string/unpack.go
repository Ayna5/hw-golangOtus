package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Trim(s, suffix string) string {
	return strings.TrimPrefix(s, suffix)
}

func Unpack(s string) (str string, err error) {
	var prevRune rune
	for _, r := range s {
		switch r {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			count, err := strconv.Atoi(string(r))
			if err != nil {
				return "", err
			}
			if unicode.IsDigit(prevRune) || strings.Index(s, string(r)) == 0 {
				return "", ErrInvalidString
			}
			if string(prevRune) != `\` && !unicode.IsDigit(prevRune) && count > 0 {
				str += strings.Repeat(string(prevRune), count-1)
			} else {
				str = Trim(str, string(prevRune))
			}
		default:
			if !unicode.IsDigit(r) {
				str += string(r)
			}
		}
		prevRune = r
	}
	return str, nil
}
