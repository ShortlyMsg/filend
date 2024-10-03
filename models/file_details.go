package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type FileDetails struct {
	FileDetailsID uuid.UUID      `gorm:"type:varchar(36);default:uuid_generate_v4();primaryKey" json:"fileDetailsId"`
	FileNames     pq.StringArray `json:"fileNames" gorm:"type:text[]"`
	FileHashes    pq.StringArray `json:"fileHashes" gorm:"type:text[]"`
	FileModelID   uuid.UUID      `gorm:"type:uuid;not null" json:"fileModelId"`
	FileModel     FileModel      `gorm:"foreignKey:FileModelID" json:"fileModel"`
}
