package services

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
)

func StartScheduler(db *gorm.DB) {
	scheduler := gocron.NewScheduler(time.Now().Location())

	// Her 1 saate silme işlemi
	scheduler.Every(1).Hours().Do(func() {
		log.Println("Eski dosyaların silinmesi için DeleteOldFiles çağrılıyor.")
		DeleteOldFiles(db)
	})

	// Zamanlayıcıyı başlat
	scheduler.StartBlocking()
}
