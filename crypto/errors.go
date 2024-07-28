package crypto

import "errors"

// HOTP
var (
	ErrInvalidDigits = errors.New("invalid digits")
	ErrInternal      = errors.New("internal error")
)

// Encryption / Decryption
var (
	ErrInitiateCipher = errors.New("failed to initiate cipher")
	ErrRandom         = errors.New("failed to generate random bytes")
)
