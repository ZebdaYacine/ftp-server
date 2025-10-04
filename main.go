package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"ftp-server/config"

	magic "github.com/ZebdaYacine/magic-bytes/magic"

	"github.com/gin-gonic/gin"
)

type UploadRequest struct {
	FileData string `json:"file_data" binding:"required"`
	FileName string `json:"file_name" binding:"required"`
}

func main() {
	env, err := config.Load()
	if err != nil {
		fmt.Println("Error loading environment variables:", err)
		return
	}

	port := env.Port
	host := env.Host
	uploadDir := env.UploadDir

	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		log.Fatalf("failed to initialize upload directory %s: %v", uploadDir, err)
	}
	router := gin.Default()

	router.POST("/upload", func(c *gin.Context) {
		var req UploadRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		dir := filepath.Join(uploadDir, req.FileName)

		if err := os.MkdirAll(filepath.Dir(dir), 0o755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to prepare directory: %v", err)})
			return
		}

		log.Println("Saving file to:", dir)

		pathURL, err := magic.SaveBase64ToFile(req.FileData, dir)
		if err != nil {
			if errors.Is(err, os.ErrPermission) {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("upload directory is not writable: %s", uploadDir)})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		fullURL := env.BaseURL + "/" + strings.TrimPrefix(*pathURL, "/var/www/ftp/")
		c.JSON(http.StatusOK, gin.H{"url": fullURL})

	})

	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Println("Server running on", addr)
	router.Run(addr)
}
