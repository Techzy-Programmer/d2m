package server

import (
	"log"

	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/internal/handler"
	"github.com/gin-gonic/gin"
)

func StartWebServer(port string) {
	router := gin.Default()
	router.Use(handler.HandleUI())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": "Wohoo! Nothing to be found here",
			"ok":      false,
		})
	})

	api := router.Group("/api")
	{
		api.POST("/deploy", handler.HandleDeployment)
		api.POST("/auth", handler.HandleAuth)
		handlePostAuthAPI(*api.Group("/mg"))
	}

	paint.Info("Serving backend at :" + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func handlePostAuthAPI(fh gin.RouterGroup) {
	fh.Use(handler.VerifySession)
	fh.GET("/meta", handler.HandleMeta)
	fh.GET("/logout", handler.HandleLogout)
	fh.GET("/get-deployments", handler.HandleGetDeployments)
}
