package helpers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
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
