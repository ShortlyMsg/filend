package controllers

import (
	"filend/services"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var files = map[string]string{}

func UploadFile(c *gin.Context) {

	// Klasör yoksa oluştur
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	otp := services.GenerateOneTimePassword()

	// Dosyayı kaydet
	filePath := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Dosya yüklenemedi"})
		return
	}

	// OTP ve dosya adını kaydet
	files[otp] = filePath

	c.JSON(http.StatusOK, gin.H{"otp": otp, "fileName": file.Filename})
}

func DownloadFile(c *gin.Context) {
	otp := c.Param("otp")

	filePath, exists := files[otp]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dosya bulunamadı"})
		return
	}

	fileName := filepath.Base(filePath)

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.File(filePath)
}
