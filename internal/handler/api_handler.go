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
			"message": "Request payload was found defective",
			"ok":      false,
		})
		return
	}

	compErr := bcrypt.CompareHashAndPassword([]byte(accessPwd), []byte(decPayload))
	if compErr != nil {
		c.JSON(401, gin.H{
			"message": "Credentials do not match",
			"ok":      false,
		})
		return
	}

	tok, tokErr := helpers.GenerateJWTToken(db.GetConfig("app.JWTSecret", ""))
	if tokErr != nil {
		c.JSON(500, gin.H{
			"message": "Internal server error (Token generation)",
			"ok":      false,
		})
		return
	}

	exp := 3600 * 24 * 7
	c.SetCookie("access_token", tok, exp, "/api", "", false, true)

	c.JSON(200, gin.H{
		"message": "Welcome to D2M Web panel experience",
		"ok":      true,
	})
}

func VerifySession(c *gin.Context) {
	cookieToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(401, gin.H{
			"message": "No active session, please authenticate",
			"ok":      false,
		})
		c.Abort()
		return
	}

	_, tokErr := helpers.VerifyJWTToken(cookieToken, db.GetConfig("app.JWTSecret", ""))
	if tokErr != nil {
		c.JSON(401, gin.H{
			"message": "Session expired, please re-authenticate",
			"ok":      false,
		})
		c.Abort()
		return
	}

	c.Next()
}

func HandleMeta(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Meta data",
		"ok":      true,
	})
}
