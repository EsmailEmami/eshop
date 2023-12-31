package models

import "github.com/google/uuid"

type ProductFileMap struct {
	ProductID uuid.UUID `gorm:"column:product_id"                   json:"productId"`
	FileID    uuid.UUID `gorm:"column:file_id"                      json:"fileId"`
	Product   Product   `gorm:"foreignKey:product_id;references:id" json:"product"`
	File      File      `gorm:"foreignKey:file_id;references:id"    json:"file"`
	Priority  int       `gorm:"column:priority"                     json:"priority"`
}

func (ProductFileMap) TableName() string {
	return "product_file_map"
}
