package utils

import (
	"encoding/base32"
)

func StringBase32ToBytes(s string) ([]byte, error) {
	base32Decoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	return base32Decoder.DecodeString(s)
}
