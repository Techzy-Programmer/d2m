package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

func generateSecureRandomString(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}
