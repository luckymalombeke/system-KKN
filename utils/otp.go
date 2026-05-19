package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateOTP6 menghasilkan kode numerik 6 digit dengan crypto/rand.
func GenerateOTP6() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
