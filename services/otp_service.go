package services

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func GenerateOneTimePassword() string {
	length := 6
	otp := make([]byte, length)
	for i := range otp {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		otp[i] = charset[randomIndex.Int64()]
	}
	return string(otp)
}
