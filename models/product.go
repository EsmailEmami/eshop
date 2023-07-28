package models

import (
	"github.com/google/uuid"
)

type Product struct {
	Model

	Name             string           `gorm:"name"                                   json:"name"`
	Code             string           `gorm:"code"                                   json:"code"`
	BrandID          uuid.UUID        `gorm:"brand_id"                               json:"brandId"`
	Brand            Brand            `gorm:"foreignKey:brand_id; references:id;"    json:"brand"`
	Files            []ProductFileMap `gorm:"foreignKey:product_id; references:id;"  json:"files"`
	ProductItems     []ProductItem    `gorm:"foreignKey:product_id; references:id;"  json:"productItems"`
	CategoryID       uuid.UUID        `gorm:"category_id"                            json:"categoryId"`
	Category         Category         `gorm:"foreignKey:category_id; references:id;" json:"category"`
	Description      string           `gorm:"description"                            json:"description"`
	ShortDescription string           `gorm:"short_description"                      json:"shortDescription"`
}

func (Product) TableName() string {
	return "product"
}
