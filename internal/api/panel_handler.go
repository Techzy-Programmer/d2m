package api

import "github.com/gin-gonic/gin"

func HandlePanel(c *gin.Context) {
	// Handle panel requests

	c.JSON(501, gin.H{
		"message": "Not implemented yet",
		"ok":      false,
	})
}
