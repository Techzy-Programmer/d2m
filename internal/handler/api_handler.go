package handler

import (
	"io"
	"strconv"

	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/vars"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type meta struct {
	WebPort string `json:"webPort"`
	TcpPort string `json:"tcpPort"`
	Uptime  string `json:"uptime"`
}

func HandleHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "d2m server is healthy",
		"uptime":  getRelativeDuration(vars.StartedAt),
		"version": vars.Version,
		"ok":      true,
	})
}

func HandleAuth(c *gin.Context) {
	accessPwd := db.GetConfig("user.AccessPwd", "")
	if accessPwd == "" {
		c.JSON(400, gin.H{
			"message": "Access password not set",
			"ok":      false,
		})
		return
	}

	pwdBytes, rdErr := io.ReadAll(c.Request.Body)
	if rdErr != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body",
			"ok":      false,
		})
		return
	}

	compErr := bcrypt.CompareHashAndPassword([]byte(accessPwd), pwdBytes)
	if compErr != nil {
		c.JSON(401, gin.H{
			"message": "Credentials do not match",
			"ok":      false,
		})
		return
	}

	tok, tokErr := generateJWTToken(db.GetConfig("app.JWTSecret", ""))
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
		"meta":    getMeta(),
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

	_, tokErr := verifyJWTToken(cookieToken, db.GetConfig("app.JWTSecret", ""))
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
		"message": "Meta data retrieved successfully",
		"meta":    getMeta(),
		"ok":      true,
	})
}

func HandleLogout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/api", "", false, true)
	c.Redirect(302, "/auth")
}

func getMeta() *meta {
	return &meta{
		WebPort: db.GetConfig("user.WebPort", ""),
		TcpPort: db.GetConfig("daemon.Port", ""),
		Uptime:  getRelativeDuration(vars.StartedAt),
	}
}

func HandleGetDeployments(c *gin.Context) {
	deployments := db.GetAllDeployments()

	c.JSON(200, gin.H{
		"message":     "Deployments fetched successfully",
		"deployments": deployments,
		"ok":          true,
	})
}

func HandleGetDeploymentDetails(c *gin.Context) {
	deployIdStr := c.Param("deployID")
	uint32DepId, err := strconv.ParseUint(deployIdStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid deployment ID",
			"ok":      false,
		})
		return
	}

	deployId := uint(uint32DepId)
	logs := db.GetAllLogsForDeployment(deployId)
	deployment, depErr := db.GetDeploymentByID(deployId)
	if depErr != nil || logs == nil {
		c.JSON(404, gin.H{
			"message": "Deployment not found",
			"ok":      false,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Deployment details retrieved",
		"meta":    deployment,
		"logs":    logs,
		"ok":      true,
	})
}
