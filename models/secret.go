package models

type Secret struct {
	Secret      []byte
	AlgoType    AlgoType
	HashType    HashType
	CountFactor uint64 // last counter for HOTP, period in second for TOTP
}
