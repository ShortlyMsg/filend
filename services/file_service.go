package services

import (
	"context"
	"filend/config"
	"filend/models"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func DeleteOldFiles(db *gorm.DB) {
	thresholdTime := time.Now().Add(-12 * time.Hour)

	var fileModels []models.FileModel
	if err := db.Where("created_at < ?", thresholdTime).Find(&fileModels).Error; err != nil {
		log.Println("Veritabanından dosya modeli alınamadı:", err)
		return
	}

	c := context.Background()

	for _, fileModel := range fileModels {

		if fileModel.DeletedAt != nil {
			continue
		}

		// Silinme tarihini güncelle
		now := time.Now()
		fileModel.DeletedAt = &now
		if err := db.Save(&fileModel).Error; err != nil {
			log.Println("Veritabanında silinme tarihi güncellenemedi:", err)
		}

		// Dosyaların hash'ine göre silinmesi
		var fileDetails []models.FileDetails
		if err := db.Where("file_model_id = ?", fileModel.FileModelID).Find(&fileDetails).Error; err != nil {
			log.Println("FileDetails bulunamadı:", err)
			continue
		}

		for _, fileDetail := range fileDetails {
			for _, hash := range fileDetail.FileHashes {
				log.Printf("Silme işlemi başlatılıyor: %s hash ile", hash)
				err := config.MinioClient.RemoveObject(c, "filend", hash, minio.RemoveObjectOptions{})
				if err != nil {
					log.Printf("MinIO'dan dosya %s silinemedi: %v", hash, err)
					continue
				}
			}
		}
	}

	log.Println("Eski dosyalar başarıyla silindi ve veritabanı güncellendi.")
}
