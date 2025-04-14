package controllers

import (
	//"context"
	"crypto/sha256"
	"encoding/hex"
	"filend/config"
	"filend/models"
	"filend/services"
	"fmt"
	"log"
	"os"
	"strconv"

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

func GenerateOTP(c *gin.Context) {
	otp := services.GenerateOneTimePassword()
	c.JSON(http.StatusOK, gin.H{"otp": otp})
}

func UpdateFileTimeByHash(fileHash string) error {

	query := `
	UPDATE file_models
	SET updated_at = NOW()
	FROM file_details
	WHERE file_details.file_model_id = file_models.file_model_id
	AND file_details.file_hash @> ? 
	AND file_models.deleted_at IS NULL;
`
	err := config.DB.Exec(query, "%"+fileHash+"%").Error
	return err
}

func UploadFile(c *gin.Context) {

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uploadedFiles := form.File["files"] // Formdan dosyaları al
	otp := c.Query("otp")
	// userSecurityCode := c.PostForm("userSecurityCode")

	// if !validateSecurityCode(userSecurityCode) {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz güvenlik kodu. Lütfen 4 haneli büyük/küçük harf ve rakam içeren bir kod girin."})
	// 	return
	// }

	fileName := c.PostForm("fileName")
	fileHash := c.PostForm("fileHash")

	chunkIndex := c.PostForm("chunkIndex")
	totalChunks := c.PostForm("totalChunks")

	var fileModel models.FileModel
	err = config.DB.Where("otp = ?", otp).First(&fileModel).Error
	if err != nil {
		fileModel = models.FileModel{
			FileModelID:      uuid.New(),
			Otp:              otp,
			UserSecurityCode: "usc0",
		}

		if err := config.DB.Create(&fileModel).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Veritabanına kaydedilemedi"})
			return
		}
	}

	for _, file := range uploadedFiles {
		// İstemciden gelecek şekilde
		uploadedFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Dosya açılamadı: " + err.Error()})
			return
		}
		defer uploadedFile.Close()

		// Chunk'ı tmp klasörüne OTP_HASH_ChunkNum şeklinde kaydediyoruz
		os.MkdirAll(fmt.Sprintf("tmp/%s", otp), os.ModePerm)
		chunkPath := fmt.Sprintf("tmp/%s/%s_chunk%s", otp, fileHash, chunkIndex)
		out, err := os.Create(chunkPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Chunk Yazılamadı: " + err.Error()})
			return
		}
		_, err = io.Copy(out, uploadedFile)
		out.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Chunk yazılırken bir hata oluştu"})
		}

		// Chunk Yüklenmesi Bitince Birleştir
		tChunks, _ := strconv.Atoi(totalChunks)
		cIndex, _ := strconv.Atoi(chunkIndex)

		if cIndex+1 == tChunks {

			mergedPath := fmt.Sprintf("tmp/%s/%s_merged", otp, fileHash)
			mergedFile, err := os.Create(mergedPath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Geçici birleştirme dosyası oluşturulamadı"})
			}
			defer mergedFile.Close()

			for i := 0; i < tChunks; i++ {
				chunkFilePath := fmt.Sprintf("tmp/%s/%s_chunk%d", otp, fileHash, i)
				chunkFile, err := os.Open(chunkFilePath)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Chunk %d okunamadı", i)})
				}
				io.Copy(mergedFile, chunkFile)
				chunkFile.Close()
				os.Remove(chunkFilePath)
			}

			// Hashi oluştur
			hash, err := GenerateFileHash(mergedFile)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Dosya Hashi oluşturulamadı"})
				return
			}

			var existingFile models.FileDetails
			err = config.DB.Where("file_details.file_hash @> ? AND file_models.deleted_at IS NULL", hash).
				Joins("JOIN file_models ON file_details.file_model_id = file_models.file_model_id").
				First(&existingFile).Error
			if err == nil {
				log.Printf("3 if err==nil içi %s", fileHash)
				// // Dosya zaten var, MinIO'ya yüklemiyoruz ama DB'ye kaydediyoruz
				// if err := UpdateFileTimeByHash(fileHash); err != nil {
				// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "UpdatedAt güncellenemedi"})
				// 	return
				// }
				// log.Printf("4 oldu gibi")
				// MinIO'ya hash ismiyle yükle
			} else {
				mergedFile.Seek(0, io.SeekStart)
				fileStat, _ := mergedFile.Stat()
				_, err = config.MinioClient.PutObject(c, "filend", hash, mergedFile, fileStat.Size(), minio.PutObjectOptions{
					ContentType: file.Header.Get("Content-Type"),
				})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "MinIO'ya yüklenemedi: " + err.Error()})
					return
				}

				fileName = file.Filename
			}
			fileDetail := models.FileDetails{
				FileDetailsID: uuid.New(),
				FileName:      fileName,
				FileHash:      fileHash,
				FileModelID:   fileModel.FileModelID,
			}

			if err := config.DB.Create(&fileDetail).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Veritabanına kaydedilemedi"})
				return
			}

			// Geçici Birleşmiş Dosyasyı Sil
			os.Remove(mergedPath)
		}

	}

	c.JSON(http.StatusOK, gin.H{"otp": otp, "fileName": fileName})
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

	if fileModel.DeletedAt != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu dosya silindiği için indirilemez"})
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
		if detail.FileHash == request.FileHash {
			fileName = detail.FileName
			break
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

	if fileModel.DeletedAt != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu dosya silindiği için indirilemez"})
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
		fileNames = append(fileNames, detail.FileName)
		fileHashes = append(fileHashes, detail.FileHash)
	}

	c.JSON(http.StatusOK, gin.H{"files": fileNames, "hashes": fileHashes})
}

func CheckFileHash(c *gin.Context) {
	var requestHash struct {
		FileHash string `json:"fileHash"`
	}

	if err := c.ShouldBindJSON(&requestHash); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek"})
		return
	}

	if len(requestHash.FileHash) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dosya hash'i boş olamaz"})
		return
	}

	var existingFile string

	err := config.DB.Table("file_details").Select("file_hashes").Where("file_details.file_hashes @> ? AND file_models.deleted_at IS NULL", requestHash.FileHash).
		Joins("JOIN file_models ON file_details.file_model_id = file_models.file_model_id").
		Scan(&existingFile).Error

	fileStatus := make(map[string]bool)

	// Eğer dosya bulunursa false, bulunamazsa true
	if err == nil && existingFile != "" {
		err := UpdateFileTimeByHash(requestHash.FileHash)
		if err != nil {
			log.Printf("Hata: Dosya zaman güncellemesi başarısız: %v", err)
		} else {
			log.Printf("Zaman güncellemesi başarıyla yapıldı: %s", requestHash.FileHash)
		}
		fileStatus[requestHash.FileHash] = false // Dosya mevcut
	} else {
		fileStatus[requestHash.FileHash] = true // Dosya mevcut değil, yüklenebilir
	}

	//log.Printf("Existing Files: %+v", existingFiles)
	c.JSON(http.StatusOK, gin.H{"fileStatus": fileStatus})
}
