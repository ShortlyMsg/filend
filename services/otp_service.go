package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOneTimePassword() string {
	num, _ := rand.Int(rand.Reader, big.NewInt(900000))
	otp := 100000 + num.Int64()
	return fmt.Sprintf("%d", otp)
}
