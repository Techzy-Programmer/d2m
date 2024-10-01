package handler

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"slices"

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

	if univ.PrivKey == nil {
		paint.Error("Private key not yet configured")
		c.JSON(500, gin.H{
			"message": "Internal server error",
			"code":    "private_key_error",
			"ok":      false,
		})
		return
	}

	decBodyBits, decryptErr := rsa.DecryptPKCS1v15(nil, univ.PrivKey, encBodyBits)
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
