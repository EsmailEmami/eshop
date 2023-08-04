package models

import "github.com/google/uuid"

type FavoriteProductItem struct {
	Model

	ProductItemID uuid.UUID    `gorm:"column:product_item_id"                   json:"productItemId"`
	ProductItem   *ProductItem `gorm:"foreignKey:product_item_id;references:id" json:"productItem"`
}

func (FavoriteProductItem) TableName() string {
	return "favorite_product_item"
}
