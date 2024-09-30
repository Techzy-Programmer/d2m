package ui

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/Techzy-Programmer/d2m/config/db"
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
		defer c.Abort()

		if strings.HasSuffix(c.Request.URL.Path, "__API-Port") {
			c.Data(200, "text/plain", []byte(db.GetConfig[string]("user.APIPort", "8080")))
			return
		}

		fileserver.ServeHTTP(c.Writer, c.Request)
	}
}

func StartUIServer(port string) {
	router := gin.Default()
	router.Use(embedReact("/", "panel/dist", web.EmbeddedFiles))

	paint.Info("Serving ui at :" + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
