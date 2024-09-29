package api

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"slices"

	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/flow"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/univ"
	"github.com/gin-gonic/gin"
)

func HandleDeployment(c *gin.Context) {
	// Following check ensures that request only gets through if coming from GitHub actions workflow
	// ToDo: Support custom CD servers and/or workflow runners
	if !slices.Contains(univ.GHActionIps, c.ClientIP()) {
		c.JSON(403, gin.H{
			"message": "You are not allowed to access this resource",
			"ok":      false,
		})
		return
	}

	var b64BodyBits []byte
	_, readErr := c.Request.Body.Read(b64BodyBits)
	if readErr != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body",
			"ok":      false,
		})
		return
	}

	var encBodyBits []byte
	_, decodeErr := base64.StdEncoding.Decode(encBodyBits, b64BodyBits)
	if decodeErr != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body",
			"ok":      false,
		})
		return
	}

	privateKey, keyErr := getPrivateKey()
	if keyErr != nil {
		paint.Error("Error getting private key: ", keyErr)
		c.JSON(500, gin.H{
			"message": "Internal server error",
			"code":    "private_key_error",
			"ok":      false,
		})
		return
	}

	decBodyBits, decryptErr := rsa.DecryptPKCS1v15(nil, privateKey, encBodyBits)
	if decryptErr != nil {
		paint.Error("Error decrypting request body: ", decryptErr)
		c.JSON(500, gin.H{
			"message": "Internal server error",
			"code":    "decryption_error",
			"ok":      false,
		})
		return
	}

	// Deserialize the decrypted request body
	var req univ.DeploymentRequest
	jsonErr := json.Unmarshal(decBodyBits, &req)
	if jsonErr != nil {
		paint.Error("Error deserializing request body: ", jsonErr)
		c.JSON(400, gin.H{
			"message": "Invalid request body",
			"ok":      false,
		})
		return
	}

	flow.StartDeployment(&req)

	c.JSON(200, gin.H{
		"message": "Deployment triggered successfully",
		"ok":      true,
	})
}

func getPrivateKey() (*rsa.PrivateKey, error) {
	// Decode the PEM block
	block, _ := pem.Decode([]byte(db.GetConfig[string]("user.PrivateKey")))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	// Parse the private key
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("error parsing private key: " + (err.Error()))
	}

	// Type assertion to convert to *rsa.PrivateKey
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}

	return rsaPrivateKey, nil
}
