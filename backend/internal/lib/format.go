package lib

import (
	"finstat/internal/apperr"
	"strings"
	"unicode"
	"unicode/utf8"
)

func FormatName(name string, length int) (string, error) {
	newName := strings.TrimSpace(name)

	if utf8.RuneCountInString(newName) < length {
		return "", apperr.ShortString
	}

	newName = strings.ToLower(newName)

	r, size := utf8.DecodeRuneInString(newName)

	return string(unicode.ToUpper(r)) + newName[size:], nil
}
