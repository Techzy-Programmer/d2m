package flow

import (
	"crypto/rand"
	"math/big"
)

func generateSecure4DigitNumber() uint {
	min := 1000
	max := 9999
	rangeVal := big.NewInt(int64(max - min + 1))

	n, err := rand.Int(rand.Reader, rangeVal)
	if err != nil {
		return 0
	}

	return uint(n.Int64() + int64(min))
}
