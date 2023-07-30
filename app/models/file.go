package models

import (
	"time"

	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
)

type FileOutPutModel struct {
	ID             *uuid.UUID      `gorm:"id"                               json:"id"`
	CreatedAt      time.Time       `gorm:"column:created_at"                json:"createdAt"`
	UpdatedAt      time.Time       `gorm:"column:updated_at"                json:"updatedAt"`
	MimeType       string          `gorm:"mime_type"                        json:"mimeType"`
	Extension      string          `gorm:"extension"                        json:"extension"`
	OriginalName   string          `gorm:"original_name"                    json:"originalName"`
	UniqueFileName string          `gorm:"unique_file_name"                 json:"uniqueFineName"`
	FileType       models.FileType `gorm:"file_type"                        json:"fileType"`
}
