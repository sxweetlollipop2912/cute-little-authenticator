package main

import (
	"little-auth/crypto/hotp"
	"little-auth/models"
	"little-auth/utils"
)

func main() {
	var (
		authSecret = models.Secret{
			AlgoType:    models.TOTP,
			CountFactor: 30,
			HashType:    models.SHA1,
		}
		otp uint32
		err error
	)
	if authSecret.Secret, err = utils.StringBase32ToBytes("6YJZ2Q3WEVVQLIBE"); err != nil {
		panic(err)
	}

	if otp, err = hotp.New(models.HashFnFromType(authSecret.HashType)).Generate(
		authSecret.Secret,
		utils.TotpCounterFromNow(authSecret.CountFactor),
		6,
	); err != nil {
		panic(err)
	}

	println(otp)
}

//package main
//
//import (
//	"fmt"
//	"github.com/99designs/keyring"
//)
//
//func main() {
//	ring, _ := keyring.Open(keyring.Config{
//		ServiceName: "example",
//	})
//
//	_ = ring.Set(keyring.Item{
//		Key:  "foo",
//		Data: []byte("secret-bar"),
//	})
//
//	i, _ := ring.Get("foo")
//
//	//_ = ring.Remove("foo")
//
//	fmt.Printf("%s", i.Data)
//}
