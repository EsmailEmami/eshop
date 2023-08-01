package models

import (
	"github.com/google/uuid"
)

type Product struct {
	Model

	Name                 string                `gorm:"column:name"                            json:"name"`
	Code                 string                `gorm:"column:code"                            json:"code"`
	BrandID              uuid.UUID             `gorm:"column:brand_id"                        json:"brandId"`
	Brand                Brand                 `gorm:"foreignKey:brand_id; references:id;"    json:"brand"`
	Files                []ProductFileMap      `gorm:"foreignKey:product_id; references:id;"  json:"files"`
	ProductItems         []ProductItem         `gorm:"foreignKey:product_id; references:id;"  json:"productItems"`
	Features             []ProductFeatureValue `gorm:"foreignKey:product_id; references:id;"  json:"features"`
	CategoryID           uuid.UUID             `gorm:"column:category_id"                     json:"categoryId"`
	Category             Category              `gorm:"foreignKey:category_id; references:id;" json:"category"`
	Description          string                `gorm:"column:description"                     json:"description"`
	ShortDescription     string                `gorm:"column:short_description"               json:"shortDescription"`
	DefaultProductItemID *uuid.UUID            `gorm:"column:default_product_item_id"         json:"defaultProductItemId"`
	Comments             []Comment             `gorm:"foreignKey:product_id;references:id"    json:"comments"`
}

func (Product) TableName() string {
	return "product"
}
