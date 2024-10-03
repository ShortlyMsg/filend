package models

import (
	"time"

	"github.com/google/uuid"
)

type FileModel struct {
	FileModelID      uuid.UUID  `gorm:"type:varchar(36);default:uuid_generate_v4();primaryKey" json:"fileModelId"`
	Otp              string     `json:"otp"`
	UserSecurityCode string     `json:"userSecurityCode"`
	CreatedAt        time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt        time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	DeletedAt        *time.Time `gorm:"index"`
}
