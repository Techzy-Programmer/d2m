package handler

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/Techzy-Programmer/d2m/web"
	"github.com/gin-gonic/gin"
)

func HandleUI() gin.HandlerFunc {
	em := web.EmbeddedFiles
	buildDirectory := "panel/dist"
	embedDir, _ := fs.Sub(em, buildDirectory)
	fileserver := http.FileServer(http.FS(embedDir))
	fileserver = http.StripPrefix("/", fileserver)

	return func(c *gin.Context) {
		path := &c.Request.URL.Path
		if strings.HasPrefix(*path, "/api") && (len(*path) == 4 || (*path)[4] == '/') {
			return
		}

		filePath := buildDirectory + *path
		_, err := em.Open(filePath)

		if err != nil {
			*path = "/"
			// If file doesn't exist in FS,
			// serve index.html for React client-side routing
		}

		fileserver.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
