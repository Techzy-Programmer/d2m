package handler

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Techzy-Programmer/d2m/config/vars"
	"github.com/golang-jwt/jwt/v5"
)

func bodyAsText(req *http.Request) string {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return ""
	}

	return string(body)
}

func getRelativeDuration(startTime int64) string {
	now := time.Now().Unix()
	duration := now - startTime

	seconds := duration
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24
	weeks := days / 7
	months := days / 30
	years := days / 365

	switch {
	case years > 0:
		return fmt.Sprintf("%d years", years)
	case months > 0:
		return fmt.Sprintf("%d months", months)
	case weeks > 0:
		return fmt.Sprintf("%d weeks", weeks)
	case days > 0:
		return fmt.Sprintf("%d days", days)
	case hours > 0:
		return fmt.Sprintf("%d hours", hours)
	case minutes > 0:
		return fmt.Sprintf("%d minutes", minutes)
	default:
		return fmt.Sprintf("%d seconds", seconds)
	}
}

func rsaDecryptWithPrivateKey(b64Ciper string) (string, error) {
	cipherBytes, b64Err := base64.StdEncoding.DecodeString(b64Ciper)
	if b64Err != nil {
		return "", errors.New("error decoding base64 string: " + b64Err.Error())
	}

	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, vars.PrivKey, cipherBytes)
	if err != nil {
		return "", errors.New("error decrypting data: " + err.Error())
	}

	return string(decryptedData), nil
}

func generateJWTToken(secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.Claims(jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "d2m-cli",
	}))

	return token.SignedString([]byte(secret))
}

func verifyJWTToken(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
