package controllers

import (
	//"context"
	"crypto/sha256"
	"encoding/hex"
	"filend/config"
	"filend/models"
	"filend/services"
	"time"

	"github.com/lib/pq"

	"io"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func validateSecurityCode(code string) bool {
	re := regexp.MustCompile(`^[A-Za-z0-9]{4}$`)
	return re.MatchString(code)
}

func GenerateFileHash(file io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func UploadFile(c *gin.Context) {

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uploadedFiles := form.File["files"] // Formdan dosyaları al
	otp := services.GenerateOneTimePassword()
	// userSecurityCode := c.PostForm("userSecurityCode")

	// if !validateSecurityCode(userSecurityCode) {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz güvenlik kodu. Lütfen 4 haneli büyük/küçük harf ve rakam içeren bir kod girin."})
	// 	return
	// }

	fileNames := c.PostFormArray("fileNames[]")
	fileHashes := c.PostFormArray("fileHashes[]")

	fileModel := models.FileModel{
		FileModelID:      uuid.New(),
		Otp:              otp,
		UserSecurityCode: "usc0",
	}

	if err := config.DB.Create(&fileModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veritabanına kaydedilemedi"})
		return
	}

	for _, file := range uploadedFiles {
		// İstemciden gelecek şekilde
		uploadedFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Dosya açılamadı: " + err.Error()})
			return
		}
		defer uploadedFile.Close()

		// Hashi oluştur
		fileHash, err := GenerateFileHash(uploadedFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Dosya Hashi oluşturulamadı"})
			return
		}

		var existingFile models.FileDetails
		err = config.DB.Where("file_details.file_hashes @> ? AND file_models.deleted_at IS NULL", pq.StringArray{fileHash}).
			Joins("JOIN file_models ON file_details.file_model_id = file_models.file_model_id").
			First(&existingFile).Error
		if err == nil {
			// Dosya zaten var, MinIO'ya yüklemiyoruz ama DB'ye kaydediyoruz
			fileNames = append(fileNames, file.Filename)
			fileHashes = append(fileHashes, fileHash)

			existingFile.FileModel.UpdatedAt = time.Now()
			if err := config.DB.Save(&existingFile.FileModel).Error; err != nil {
				return
			}
		} else {
			// MinIO'ya hash ismiyle yükle
			uploadedFile.Seek(0, io.SeekStart)
			_, err = config.MinioClient.PutObject(c, "filend", fileHash, uploadedFile, file.Size, minio.PutObjectOptions{
				ContentType: file.Header.Get("Content-Type"),
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "MinIO'ya yüklenemedi: " + err.Error()})
				return
			}

			fileNames = append(fileNames, file.Filename)
			fileHashes = append(fileHashes, fileHash)
		}
	}
	fileDetail := models.FileDetails{
		FileDetailsID: uuid.New(),
		FileNames:     pq.StringArray(fileNames),
		FileHashes:    pq.StringArray(fileHashes),
		FileModelID:   fileModel.FileModelID,
	}

	if err := config.DB.Create(&fileDetail).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veritabanına kaydedilemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"otp": otp, "fileNames": fileNames})
}

func DownloadFile(c *gin.Context) {

	var request struct {
		Otp      string `json:"otp"`
		FileHash string `json:"fileHash"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek"})
		return
	}
	// userSecurityCode := c.Query("userSecurityCode")

	// if !validateSecurityCode(userSecurityCode) {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz güvenlik kodu. Lütfen 4 haneli büyük/küçük harf ve rakam içeren bir kod girin."})
	// 	return
	// }

	var fileModel models.FileModel
	if err := config.DB.Where("otp = ?", request.Otp).First(&fileModel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dosya bulunamadı"})
		return
	}

	var fileDetails []models.FileDetails
	if err := config.DB.Where("file_model_id = ?", fileModel.FileModelID).Find(&fileDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "FileDetails Getirilemedi"})
		return
	}

	if request.FileHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dosya hashi belirtilmedi"})
		return
	}

	var fileName string
	for _, detail := range fileDetails {
		for i, fileHash := range detail.FileHashes {
			if fileHash == request.FileHash {
				fileName = detail.FileNames[i]
				break
			}
		}
	}

	// MinIO'dan dosyayı getir ve indir
	fileObject, err := config.MinioClient.GetObject(c, "filend", request.FileHash, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Dosya getirilemedi: " + err.Error()})
		return
	}
	defer fileObject.Close()

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	io.Copy(c.Writer, fileObject)
}

func GetAllFiles(c *gin.Context) {
	var request struct {
		Otp string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek"})
		return
	}

	otp := request.Otp

	var fileModel models.FileModel
	if err := config.DB.Where("otp = ?", otp).First(&fileModel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dosyalar bulunamadı"})
		return
	}

	var fileDetails []models.FileDetails
	if err := config.DB.Where("file_model_id = ?", fileModel.FileModelID).Find(&fileDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "FileDetails getirilemedi"})
		return
	}

	var fileNames []string
	var fileHashes []string
	for _, detail := range fileDetails {
		fileNames = append(fileNames, detail.FileNames...)
		fileHashes = append(fileHashes, detail.FileHashes...)
	}

	c.JSON(http.StatusOK, gin.H{"files": fileNames, "hashes": fileHashes})
}

func CheckFileHash(c *gin.Context) {
	var requestHashes struct {
		FileHashes []string `json:"fileHashes"`
	}

	if err := c.ShouldBindJSON(&requestHashes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek"})
		return
	}

	if len(requestHashes.FileHashes) > 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "En fazla 20 dosya kontrol edilebilir"})
		return
	}

	fileStatus := make(map[string]bool)

	for _, fileHash := range requestHashes.FileHashes {
		var existingFiles []models.FileDetails

		err := config.DB.Where("file_details.file_hashes @> ? AND file_models.deleted_at IS NULL", pq.StringArray{fileHash}).
			Joins("JOIN file_models ON file_details.file_model_id = file_models.file_model_id").
			Find(&existingFiles).Error

		// Eğer dosya bulunursa false, bulunamazsa true
		if err == nil && len(existingFiles) > 0 {
			fileStatus[fileHash] = false // Dosya mevcut
		} else {
			fileStatus[fileHash] = true // Dosya mevcut değil, yüklenebilir
		}
	}

	//log.Printf("Existing Files: %+v", existingFiles)
	c.JSON(http.StatusOK, gin.H{"fileStatus": fileStatus})
}
