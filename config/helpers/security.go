package helpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"time"

	"github.com/Techzy-Programmer/d2m/config/vars"
	"github.com/golang-jwt/jwt/v5"
)

func GetPrivateKey(privText string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privText))
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	var privateKey interface{}
	var err error

	switch block.Type {
	case "PRIVATE KEY": // PKCS#8
		privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	case "RSA PRIVATE KEY": // PKCS#1
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	default:
		return nil, errors.New("unsupported PEM block type")
	}

	if err != nil {
		return nil, errors.New("error parsing private key: " + err.Error())
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}

	return rsaPrivateKey, nil
}

func RSADecryptWithPrivateKey(b64Ciper string) (string, error) {
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

func GenerateSecureRandomString(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

func GenerateJWTToken(secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.Claims(jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "d2m-cli",
	}))

	return token.SignedString([]byte(secret))
}

func VerifyJWTToken(tokenString, secret string) (jwt.MapClaims, error) {
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
