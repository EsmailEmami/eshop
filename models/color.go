package models

type Color struct {
	Model

	Name         string        `gorm:"column:name"                       json:"name"`
	Code         string        `gorm:"column:code"                       json:"code"`
	ColorHex     string        `gorm:"column:color_hex"                  json:"colorHex"`
	ProductItems []ProductItem `gorm:"foreignKey:color_id;references:id" json:"products"`
}

func (Color) TableName() string {
	return "color"
}
