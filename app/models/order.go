package models

type OrderOutPutModel struct {
	Price float64                `gorm:"-" json:"price"`
	Items []OrderItemOutPutModel `gorm:"-" json:"items"`
}
