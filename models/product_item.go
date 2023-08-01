package models

import "github.com/google/uuid"

type ProductItem struct {
	Model

	Price          float64       `gorm:"column:price"                             json:"price"`
	Status         ProductStatus `gorm:"column:status"                            json:"status"`
	ColorID        uuid.UUID     `gorm:"column:color_id"                          json:"colorId"`
	Color          Color         `gorm:"foreignKey:color_id; references:id;"      json:"color"`
	ProductID      uuid.UUID     `gorm:"column:product_id"                        json:"productId"`
	Product        *Product      `gorm:"foreignKey:product_id; references:id;"    json:"product"`
	Quantity       int           `gorm:"column:quantity"                          json:"quantity"`
	BoughtQuantity int           `gorm:"column:bought_quantity"                   json:"boughtQuantity"`
	OrderItems     []OrderItem   `gorm:"foreignKey:product_item_id;references:id" json:"items"`
}

func (ProductItem) TableName() string {
	return "product_item"
}

type ProductStatus int

const (
	ProductStatusPublish ProductStatus = iota
	ProductStatusInActive
)

func (ps ProductStatus) String() string {
	switch ps {
	case ProductStatusPublish:
		return "Publish"
	case ProductStatusInActive:
		return "InActive"
	default:
		return "unknown"
	}
}
