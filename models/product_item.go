package models

import "github.com/google/uuid"

type ProductItem struct {
	Model

	Price     float64       `gorm:"price"                                 json:"price"`
	Status    ProductStatus `gorm:"status"                                json:"status"`
	ColorID   uuid.UUID     `gorm:"color_id"                              json:"colorId"`
	Color     Color         `gorm:"foreignKey:color_id; references:id;"   json:"color"`
	ProductID uuid.UUID     `gorm:"product_id"                            json:"productId"`
	Product   Product       `gorm:"foreignKey:product_id; references:id;" json:"product"`
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
