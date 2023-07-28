package models

type Brand struct {
	Model

	Name     string    `gorm:"name"                              json:"name"`
	Code     string    `gorm:"code"                              json:"code"`
	Products []Product `gorm:"foreignKey:brand_id;references:id" json:"products"`
}
