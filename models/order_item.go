package models

import "github.com/google/uuid"

type OrderItem struct {
	Model

	OrderID       uuid.UUID    `gorm:"column:order_id"                          json:"orderId"`
	Order         *Order       `gorm:"foreignKey:order_id;references:id;"       json:"order"`
	ProductItemID uuid.UUID    `gorm:"column:product_item_id"                   json:"productItemId"`
	ProductItem   *ProductItem `gorm:"foreignKey:product_item_id;references:id" json:"productItem"`
	Quantity      int          `gorm:"column:quantity"                          json:"quantity"`
	Price         float64      `gorm:"column:price"                             json:"price"`
}

func (OrderItem) TableName() string {
	return "order_item"
}
