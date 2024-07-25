package models

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

type AlgoType string

const (
	TOTP AlgoType = "totp"
	HOTP AlgoType = "hotp"
)

type HashType string

const (
	SHA1   HashType = "sha1"
	SHA256 HashType = "sha256"
	SHA512 HashType = "sha512"
)

var mapHashToFn = map[HashType]func() hash.Hash{
	SHA1:   sha1.New,
	SHA256: sha256.New,
	SHA512: sha512.New,
}

func HashFnFromType(hashType HashType) func() hash.Hash {
	return mapHashToFn[hashType]
}
