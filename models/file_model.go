package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type FileModel struct {
	gorm.Model

	ID               uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Otp              string         `json:"otp"`
	UserSecurityCode string         `json:"userSecurityCode"`
	FileNames        pq.StringArray `json:"fileNames" gorm:"type:text[]"`
}
