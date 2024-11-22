package handler

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/gin-gonic/gin"
)

var tempFilePrefix = "d2m-upload-"

// init cleans up any temp files that may have been left behind
func init() {
	tempDir := os.TempDir()
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), tempFilePrefix) {
			filePath := filepath.Join(tempDir, file.Name())

			err := os.Remove(filePath)
			if err != nil {
				paint.NoticeF("Failed to delete tmp file %s: %v\n", filePath, err)
				continue
			}
		}
	}
}

func HandleUpload(c *gin.Context) {
	tempFile, err := os.CreateTemp("", tempFilePrefix+"*.zip")
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to create temp file",
			"error":   err.Error(),
			"ok":      false,
		})
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to save uploaded file",
			"ok":      false,
		})
		return
	}

	time.AfterFunc(2*time.Minute, func() {
		os.Remove(tempFile.Name())
	})

	c.JSON(200, gin.H{
		"message":  "File uploaded successfully",
		"tempFile": tempFile.Name(),
		"ok":       true,
	})
}
