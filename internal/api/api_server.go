package api

import (
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/gin-gonic/gin"
)

func StartAPIServer(port string) {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
			"ok":      true,
		})
	})

	router.Run(":" + port)
	paint.Info("API server started at :" + port)
}
