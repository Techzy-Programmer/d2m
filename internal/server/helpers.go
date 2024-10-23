package server

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/Techzy-Programmer/d2m/config/vars"
	"github.com/gin-gonic/gin"
)

func decryptionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		encryptedKeyB64 := c.GetHeader("X-Encryption-Key")
		if encryptedKeyB64 == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Request body is encrypted but missing encryption key",
				"ok":      false,
			})
			return
		}

		encryptedKey, err := base64.StdEncoding.DecodeString(encryptedKeyB64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Invalid encryption key format",
				"ok":      false,
			})
			return
		}

		if vars.PrivKey == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
				"code":    "private_key_error",
				"ok":      false,
			})
			return
		}

		aesKey, err := rsa.DecryptPKCS1v15(nil, vars.PrivKey, encryptedKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Access forbidden",
				"ok":      false,
			})
			return
		}

		encryptedBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request body",
				"ok":      false,
			})
			return
		}

		encryptedBodyBytes, err := base64.StdEncoding.DecodeString(string(encryptedBody))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Invalid encrypted body format",
				"ok":      false,
			})
			return
		}

		block, err := aes.NewCipher(aesKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
				"code":    "aes_cipher_error",
				"ok":      false,
			})
			return
		}

		if len(encryptedBodyBytes) < aes.BlockSize {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Too short encrypted body",
				"ok":      false,
			})
			return
		}

		iv := encryptedBodyBytes[:aes.BlockSize]
		encryptedBodyBytes = encryptedBodyBytes[aes.BlockSize:]

		mode := cipher.NewCBCDecrypter(block, iv)
		decryptedBody := make([]byte, len(encryptedBodyBytes))
		mode.CryptBlocks(decryptedBody, encryptedBodyBytes)

		padding := int(decryptedBody[len(decryptedBody)-1])
		decryptedBody = decryptedBody[:len(decryptedBody)-padding]

		c.Request.Body = io.NopCloser(bytes.NewBuffer(decryptedBody))
		c.Request.ContentLength = int64(len(decryptedBody))

		c.Next()
	}
}
