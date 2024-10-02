package helpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"

	"github.com/Techzy-Programmer/d2m/config/vars"
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
