package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFotoAset(c *gin.Context) {
	file, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File foto wajib diunggah"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hanya format .jpg, .jpeg, dan .png yang diizinkan"})
		return
	}

	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, 0755)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat folder penyimpanan di server"})
			return
		}
	}

	filenameBaru := uuid.New().String() + ext
	pathTujuan := filepath.Join(uploadDir, filenameBaru)

	if err := c.SaveUploadedFile(file, pathTujuan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file ke filesystem"})
		return
	}

	fileURL := fmt.Sprintf("http://127.0.0.1:8080/uploads/%s", filenameBaru)
	c.JSON(http.StatusOK, gin.H{
		"message":  "File foto berhasil disimpan di filesystem",
		"url":      fileURL,
		"filename": filenameBaru,
	})
}