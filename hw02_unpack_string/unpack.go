package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var str string
	r := []rune(s)
	if len(r) > 0 && (unicode.IsDigit(r[0]) || string(r[0]) == `\` || string(r[len(r)-1]) == `\`) {
		return "", ErrInvalidString
	}
	for i := range r {
		if i == len(r)-1 {
			if string(r[i-1]) == `\` && unicode.IsDigit(r[i]) {
				str += string(r[i])
				return str, nil
			}
			if unicode.IsDigit(r[i]) {
				return str, nil
			}
			str += string(r[i])
			return str, nil
		}
		if unicode.IsDigit(r[i]) && unicode.IsDigit(r[i+1]) {
			return "", ErrInvalidString
		}
		if (unicode.IsLetter(r[i]) && (unicode.IsLetter(r[i+1]))) || (unicode.IsLetter(r[i]) && string(r[i+1]) == `\`) || (unicode.IsDigit(r[i]) && string(r[i+1]) == `\`) {
			str += string(r[i])
		}
		if unicode.IsLetter(r[i]) && unicode.IsDigit(r[i+1]) {
			count, err := strconv.Atoi(string(r[i+1]))
			if err != nil {
				return "", err
			}
			str += strings.Repeat(string(r[i]), count)
		}
		if string(r[i]) == `\` && string(r[i+1]) == `\` && string(r[i+2]) == `\` && unicode.IsDigit(r[i+3]) {
			str += string(r[i+1]) + string(r[i+3])
			if i+3 == len(r)-1 {
				return str, nil
			}
		}
		if i+2 <= len(r)-1 && string(r[i]) == `\` && (unicode.IsDigit(r[i+1]) || unicode.IsLetter(r[i+1]) || string(r[i+1]) == `\`) && unicode.IsDigit(r[i+2]) {
			count, err := strconv.Atoi(string(r[i+2]))
			if err != nil {
				return "", err
			}
			str += strings.Repeat(string(r[i+1]), count)
			if i+2 == len(r)-1 {
				return str, nil
			}
		}
	}
	return str, nil
}
