package main

import (
	"fmt"
	"net/http"
	"strings"

	"ftp-server/config"

	magic "github.com/ZebdaYacine/magic-bytes/magic"

	"github.com/gin-gonic/gin"
)

type UploadRequest struct {
	FileData string `json:"file_data" binding:"required"`
	Folder   string `json:"folder" binding:"required"`
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
	router := gin.Default()

	router.POST("/upload", func(c *gin.Context) {
		var req UploadRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		dir := uploadDir + "/" + req.Folder

		path_url, err := magic.SaveBase64ToFile(req.FileData, dir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fullURL := env.BaseURL + "/" + strings.TrimPrefix(*path_url, "/var/www/ftp/")
		c.JSON(http.StatusOK, gin.H{"url": fullURL, "path": fullURL})

	})

	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Println("Server running on", addr)
	router.Run(addr)
}
