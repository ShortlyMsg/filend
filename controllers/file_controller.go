package controllers

import (
	"filend/config"
	"filend/models"
	"filend/services"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func validateSecurityCode(code string) bool {
	re := regexp.MustCompile(`^[A-Za-z0-9]{4}$`)
	return re.MatchString(code)
}

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
	userSecurityCode := c.PostForm("userSecurityCode")

	if !validateSecurityCode(userSecurityCode) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz güvenlik kodu. Lütfen 4 haneli büyük/küçük harf ve rakam içeren bir kod girin."})
		return
	}

	var fileNames []string

	for _, file := range uploadedFiles {
		filePath := filepath.Join("uploads", file.Filename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Dosya yüklenemedi"})
			return
		}

		fileNames = append(fileNames, file.Filename) // Dosya adını dilime ekle
	}

	fileModel := models.FileModel{
		ID:               uuid.New(),
		Otp:              otp,
		UserSecurityCode: userSecurityCode,
		FileNames:        fileNames,
	}
	if err := config.DB.Create(&fileModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veritabanına kaydedilemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"otp": otp, "fileNames": fileNames})
}

func DownloadFile(c *gin.Context) {
	otp := c.Param("otp")
	userSecurityCode := c.Query("userSecurityCode")

	if !validateSecurityCode(userSecurityCode) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz güvenlik kodu. Lütfen 4 haneli büyük/küçük harf ve rakam içeren bir kod girin."})
		return
	}

	var fileModel models.FileModel
	if err := config.DB.Where("otp = ? AND user_security_code = ?", otp, userSecurityCode).First(&fileModel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dosya bulunamadı veya güvenlik kodu yanlış"})
		return
	}

	requestedFile := c.Query("fileName")
	if requestedFile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dosya adı belirtilmedi"})
		return
	}

	var filePath string
	for _, fileName := range fileModel.FileNames {
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

	var fileModel models.FileModel
	if err := config.DB.Where("otp = ?", otp).First(&fileModel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dosyalar bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": fileModel.FileNames})
}
