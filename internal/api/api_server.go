package api

import (
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/gin-gonic/gin"
)

func StartAPIServer(port string) {
	router := gin.Default()

	router.POST("/deploy", HandleDeployment)
	router.POST("/panel", HandlePanel)

	router.Run(":" + port)
	paint.Info("API server started at :" + port)
}
