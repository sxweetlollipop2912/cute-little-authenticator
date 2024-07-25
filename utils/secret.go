package utils

import "strings"

func NormalizeSecret(secret string) string {
	return strings.ToUpper(strings.TrimSpace(secret))
}
