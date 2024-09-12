package ui

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/web"
	"github.com/gin-gonic/gin"
)

func embedReact(urlPrefix, buildDirectory string, em embed.FS) gin.HandlerFunc {
	embedDir, _ := fs.Sub(em, buildDirectory)
	fileserver := http.FileServer(http.FS(embedDir))

	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}

	return func(c *gin.Context) {
		fileserver.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

func StartUIServer(port string) {
	r := gin.Default()
	r.Use(embedReact("/", "panel/dist", web.EmbeddedFiles))

	paint.Info("Serving ui at :" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
