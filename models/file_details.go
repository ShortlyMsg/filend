package models

import (
	"github.com/google/uuid"
)

type FileDetails struct {
	FileDetailsID uuid.UUID `gorm:"type:varchar(36);default:uuid_generate_v4();primaryKey" json:"fileDetailsId"`
	FileName      string    `json:"fileName"`
	FileHash      string    `json:"fileHash"`
	FileSize      int64     `json:"fileSize"`
	IsUploaded    bool      `gorm:"default:false" json:"isUploaded"`
	FileModelID   uuid.UUID `gorm:"type:uuid;not null" json:"fileModelId"`
	FileModel     FileModel `gorm:"foreignKey:FileModelID" json:"fileModel"`
}
