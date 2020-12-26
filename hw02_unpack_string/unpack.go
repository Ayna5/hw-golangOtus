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
	for i, r := range s {
		err := validateErr(s, i)
		if err != nil {
			return "", err
		}
		switch string(r) {
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			count, err := strconv.Atoi(string(r))
			if err != nil {
				return "", err
			}
			if string(s[i-1]) != `\` && !unicode.IsDigit(rune(s[i-1])) {
				str += strings.Repeat(string(s[i-1]), count)
			}
		case `\`:
			str1, err := checkSlash(str, s, i)
			if err != nil {
				return "", err
			}
			str = str1
			return str, nil
		default:
			str = checkDefault(str, s, i)
		}
	}
	return str, nil
}

func validateErr(s string, i int) error {
	if len(s) > 0 && (unicode.IsDigit(rune(s[0])) || string(s[0]) == `\` || string(s[len(s)-1]) == `\`) {
		return ErrInvalidString
	}
	if i+1 <= len(s)-1 && unicode.IsDigit(rune(s[i])) && unicode.IsDigit(rune(s[i+1])) {
		if string(s[i-1]) != `\` {
			return ErrInvalidString
		}
	}
	return nil
}

func checkSlash(str string, s string, i int) (string, error) {
	if string(s[i]) == `\` && string(s[i+1]) == `\` && string(s[i+2]) == `\` && unicode.IsDigit(rune(s[i+3])) {
		str += string(s[i+1]) + string(s[i+3])
		if i+3 == len(s)-1 {
			return str, nil
		}
	}
	if i+2 <= len(s)-1 && string(s[i]) == `\` && (unicode.IsDigit(rune(s[i+1])) || unicode.IsLetter(rune(s[i+1])) || string(s[i+1]) == `\`) && unicode.IsDigit(rune(s[i+2])) {
		count, err := strconv.Atoi(string(s[i+2]))
		if err != nil {
			return "", err
		}
		str += strings.Repeat(string(s[i+1]), count)
	}
	if string(s[i]) == `\` && unicode.IsDigit(rune(s[i+1])) && string(s[i+2]) == `\` {
		str += string(s[i+1]) + string(s[i+3])
		if i+3 == len(s)-1 {
			return str, nil
		}
	}
	return str, nil
}

func checkDefault(str string, s string, i int) string {
	if (i+1 <= len(s)-1 && (unicode.IsLetter(rune(s[i])) && (unicode.IsLetter(rune(s[i+1]))) || (unicode.IsLetter(rune(s[i])) && string(s[i+1]) == `\`) || (unicode.IsDigit(rune(s[i])) && string(s[i+1]) == `\`))) || i == len(s)-1 {
		str += string(s[i])
	}
	return str
}
