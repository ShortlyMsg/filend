package models

import (
	"time"

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
	FileHashes       pq.StringArray `json:"fileHashes" gorm:"type:text[]"`
	CreatedAt        time.Time      `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt        time.Time      `gorm:"type:timestamptz;not null;default:now()"`
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}
