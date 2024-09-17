package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileModel struct {
	FileModelID      uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"fileModelId"`
	Otp              string         `json:"otp"`
	UserSecurityCode string         `json:"userSecurityCode"`
	CreatedAt        time.Time      `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt        time.Time      `gorm:"type:timestamptz;not null;default:now()"`
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}
