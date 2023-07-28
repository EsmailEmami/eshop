package models

import "github.com/google/uuid"

type OrderItem struct {
	Model

	OrderID       uuid.UUID   `gorm:"order_id"                                 json:"orderId"`
	Order         Order       `gorm:"foreignKey:order_id;references:id;"       json:"order"`
	ProductItemID uuid.UUID   `gorm:"product_item_id"                          json:"productItemId"`
	ProductItem   ProductItem `gorm:"foreignKey:product_item_id;references:id" json:"productItem"`
	Quantity      int         `gorm:"quantity"                                 json:"quantity"`
	Price         float64     `gorm:"price"                                    json:"price"`
}

func (OrderItem) TableName() string {
	return "order_item"
}
