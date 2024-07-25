package crypto

type OTPGenerator interface {
	Generate(secret []byte, counter uint64, digits uint) (uint32, error)
}
