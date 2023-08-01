package models

import (
	"time"

	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
)

type FileOutPutModel struct {
	ID             *uuid.UUID      `gorm:"column:id"               json:"id"`
	CreatedAt      time.Time       `gorm:"column:created_at"       json:"createdAt"`
	UpdatedAt      time.Time       `gorm:"column:updated_at"       json:"updatedAt"`
	MimeType       string          `gorm:"column:mime_type"        json:"mimeType"`
	Extension      string          `gorm:"column:extension"        json:"extension"`
	OriginalName   string          `gorm:"column:original_name"    json:"originalName"`
	UniqueFileName string          `gorm:"column:unique_file_name" json:"uniqueFineName"`
	FileType       models.FileType `gorm:"column:file_type"        json:"fileType"`
}
