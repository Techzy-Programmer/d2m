package handler

import (
	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/helpers"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HandleAuth(c *gin.Context) {
	accessPwd := db.GetConfig("user.AccessPwd", "")
	if accessPwd == "" {
		c.JSON(400, gin.H{
			"message": "Access password not set",
			"ok":      false,
		})
		return
	}

	b64Payload := helpers.BodyAsText(c.Request)
	decPayload, decErr := helpers.RSADecryptWithPrivateKey(b64Payload)
	if decErr != nil {
		c.JSON(400, gin.H{
			"message": "Bad payload",
			"ok":      false,
		})
		return
	}

	compErr := bcrypt.CompareHashAndPassword([]byte(accessPwd), []byte(decPayload))
	if compErr != nil {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
			"ok":      false,
		})
		return
	}

	// ToDo: Generate a JWT token and send it back
	c.JSON(200, gin.H{
		"message": "Authorized",
		"ok":      true,
	})
}
