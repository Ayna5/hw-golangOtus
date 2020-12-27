package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (str string, err error) {
	var prevRune rune
	if len(s) > 0 && unicode.IsDigit([]rune(s)[0]) {
		return "", ErrInvalidString
	}
	for _, r := range s {
		switch r {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			count, err := strconv.Atoi(string(r))
			if err != nil {
				return "", err
			}
			if unicode.IsDigit(prevRune) {
				return "", ErrInvalidString
			}
			if !unicode.IsDigit(prevRune) && count > 0 {
				str += strings.Repeat(string(prevRune), count-1)
			} else {
				str = strings.TrimSuffix(str, string(prevRune))
			}
		default:
			str += string(r)
		}
		prevRune = r
	}
	return str, nil
}
