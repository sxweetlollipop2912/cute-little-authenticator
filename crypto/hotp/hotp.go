package hotp

import (
	"crypto/hmac"
	"encoding/binary"
	zlog "github.com/rs/zerolog/log"
	"hash"
	"little-auth/crypto"
)

type HOTP struct {
	hmac func(key []byte) hash.Hash
}

func New(h func() hash.Hash) *HOTP {
	return &HOTP{hmac: func(key []byte) hash.Hash {
		return hmac.New(h, key)
	}}
}

// Generate follows RFC-4226 https://datatracker.ietf.org/doc/html/rfc4226#section-5
func (h *HOTP) Generate(secret []byte, counter uint64, digits uint) (uint32, error) {
	if digits < 6 || digits > 8 {
		return 0, crypto.ErrInvalidDigits
	}

	// Turn counter into an 8-byte array
	counterBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(counterBytes[:], counter)

	// Generate a 20-byte MAC value
	hmacWithKeyHashFunc := h.hmac(secret)
	_, err := hmacWithKeyHashFunc.Write(counterBytes[:])
	if err != nil {
		zlog.Err(err).Msg("failed to write counter bytes")
		return 0, crypto.ErrInternal
	}
	macValue := hmacWithKeyHashFunc.Sum(nil)

	// Truncate
	if len(macValue) != 20 {
		zlog.Error().Msgf("mac value length is not 20, but %d", len(macValue))
		return 0, crypto.ErrInternal
	}
	offset := macValue[19] & 0xf
	truncatedMacValue := binary.BigEndian.Uint32(macValue[offset:]) & 0x7fffffff

	// Generate the OTP with the given number of digits
	pow10 := uint32(1)
	for i := uint(0); i < digits; i++ {
		pow10 *= uint32(10)
	}
	otp := truncatedMacValue % pow10

	return otp, nil
}

func (h *HOTP) SecurityLevel(digits uint) uint {
	// TODO https://datatracker.ietf.org/doc/html/rfc4226#section-6
	return 0
}
