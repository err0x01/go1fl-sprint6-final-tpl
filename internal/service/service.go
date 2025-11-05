package service

import (
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func AutoDetectAndConvert(input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", nil
	}
	if isMorseCode(input) {
		return morse.ToText(input), nil
	} else {
		return morse.ToMorse(input), nil
	}
}
func isMorseCode(input string) bool {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return false
	}
	for _, char := range trimmed {
		if char != '.' && char != '-' && char != ' ' && char != '/' {
			return false
		}
	}
	hasDotsOrDashes := false
	for _, char := range trimmed {
		if char == '.' || char == '-' {
			hasDotsOrDashes = true
			break
		}
	}
	return hasDotsOrDashes
}
