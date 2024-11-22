package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/internal/handler"
	"github.com/gin-gonic/gin"
)

var StopWebServer = make(chan bool)

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
		api.Use(decryptionMiddleware())
		api.POST("/deploy", handler.HandleDeployment)
		api.PUT("/upload", handler.HandleUpload)
		api.POST("/auth", handler.HandleAuth)
		handlePostAuthAPI(*api.Group("/mg"))
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	paint.Info("Serving backend at :" + port)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			paint.Warn("Stopping web server...")
		}
	}()

	<-StopWebServer
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to stop web server: %v", err)
	}
}

func handlePostAuthAPI(fh gin.RouterGroup) {
	fh.Use(handler.VerifySession)
	fh.GET("/meta", handler.HandleMeta)
	fh.GET("/logout", handler.HandleLogout)
	fh.GET("/get-deployments", handler.HandleGetDeployments)
	fh.GET("/deployment/:deployID", handler.HandleGetDeploymentDetails)
}
