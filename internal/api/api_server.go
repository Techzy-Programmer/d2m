package api

import (
	"log"

	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/gin-gonic/gin"
)

func StartAPIServer(port string) {
	router := gin.Default()

	router.POST("/deploy", HandleDeployment)
	router.POST("/panel", HandlePanel)

	router.GET("/", func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": "Wohoo! Nothing to be found here",
			"ok":      false,
		})
	})

	paint.Info("Serving backend at :" + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
