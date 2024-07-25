package main

import (
	"crypto/sha1"
	"little-auth/crypto/hotp"
	"little-auth/utils"
	"strings"
	"time"
)

func main() {
	var (
		// mutually agreed
		secret        = "6YJZ2Q3WEVVQLIBE"
		T0     uint64 = 0
		X      uint64 = 30
		hashFn        = sha1.New

		//RFC-6238
		counter = (uint64(time.Now().UnixMilli()/1000) - T0) / X
	)

	secret = strings.ToUpper(strings.TrimSpace(secret))
	secretBytes, err := utils.StringBase32ToBytes(secret)
	otp, err := hotp.New(hashFn).Generate(secretBytes, counter, 6)
	if err != nil {
		panic(err)
	}
	println(otp)
}
