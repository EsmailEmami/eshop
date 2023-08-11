package models

type RevenueByCategory struct {
	Price        float64 `gorm:"price" json:"price"`
	CategoryName string  `gorm:"category_name" json:"categoryName"`
	Percentage   float32 `gorm:"percentage" json:"percentage"`
}
