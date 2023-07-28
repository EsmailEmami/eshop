package models

import "github.com/google/uuid"

type Brand struct {
	Model

	Name     string    `gorm:"name"                              json:"name"`
	Code     string    `gorm:"code"                              json:"code"`
	Products []Product `gorm:"foreignKey:brand_id;references:id" json:"products"`
	FileID   uuid.UUID `gorm:"file_id"                           json:"fileId"`
	File     File      `gorm:"foreignKey:file_id;references:id;" json:"file"`
}

func (Brand) TableName() string {
	return "brand"
}
