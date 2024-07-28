package aesgcm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	zlog "github.com/rs/zerolog/log"
	"io"
	"little-auth/crypto"
)

type Mode int

const (
	AES128 Mode = 16
	AES192      = 24
	AES256      = 32
)

type AESGCM struct {
	mode Mode
}

func New(mode Mode) *AESGCM {
	return &AESGCM{mode: mode}
}

func (a *AESGCM) Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	aesgcm, err := a.newCipher(key)
	if err != nil {
		zlog.Err(err).Msg("failed to create new cipher")
		return nil, crypto.ErrInitiateCipher
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		zlog.Err(err).Msg("failed to create nonce")
		return nil, crypto.ErrRandom
	}

	// add nonce as a prefix to the encrypted data
	encrypted := aesgcm.Seal(nonce, nonce, plaintext, nil)
	return encrypted, nil
}

func (a *AESGCM) Decrypt(encrypted []byte, key []byte) ([]byte, error) {
	aesgcm, err := a.newCipher(key)
	if err != nil {
		zlog.Err(err).Msg("failed to create new cipher")
		return nil, crypto.ErrInitiateCipher
	}

	// Extract the nonce from the encrypted data
	nonceSize := aesgcm.NonceSize()
	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]

	// Decrypt the data
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		zlog.Err(err).Msg("failed to decrypt data")
		return nil, crypto.ErrInternal
	}
	return plaintext, nil
}

func (a *AESGCM) newCipher(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(block)
}

func (a *AESGCM) KeyLength() int {
	return int(a.mode)
}
