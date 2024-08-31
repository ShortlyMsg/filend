package controllers

import (
	"filend/services"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var files = map[string][]string{}

func UploadFile(c *gin.Context) {

	// Klasör yoksa oluştur
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	form, err := c.MultipartForm() // Çoklu dosya
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uploadedFiles := form.File["files"] // Formdan dosyaları al
	otp := services.GenerateOneTimePassword()

	var fileNames []string

	for _, file := range uploadedFiles {
		filePath := filepath.Join("uploads", file.Filename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Dosya yüklenemedi"})
			return
		}

		fileNames = append(fileNames, file.Filename) // Dosya adını dilime ekle
	}

	files[otp] = fileNames // OTP'ye bağlı dosya adlarını kaydet

	c.JSON(http.StatusOK, gin.H{"otp": otp, "fileNames": fileNames})
}

func DownloadFile(c *gin.Context) {
	otp := c.Param("otp")

	fileNames, exists := files[otp]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dosya bulunamadı"})
		return
	}

	requestedFile := c.Query("fileName")
	if requestedFile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dosya adı belirtilmedi"})
		return
	}

	var filePath string
	for _, fileName := range fileNames {
		if fileName == requestedFile {
			filePath = filepath.Join("uploads", fileName)
			break
		}
	}

	c.Header("Content-Disposition", "attachment; filename="+requestedFile)
	c.File(filePath)
}

func GetAllFiles(c *gin.Context) {
	otp := c.Param("otp")

	fileNames, exists := files[otp]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dosyalar bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": fileNames})
}
