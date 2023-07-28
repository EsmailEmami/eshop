package models

type Color struct {
	Model

	Name         string        `gorm:"name"                              json:"name"`
	Code         string        `gorm:"code"                              json:"code"`
	ColorHex     string        `gorm:"color_hex"                         json:"colorHex"`
	ProductItems []ProductItem `gorm:"foreignKey:color_id;references:id" json:"products"`
}

func (Color) TableName() string {
	return "color"
}
