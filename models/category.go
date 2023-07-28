package models

type Category struct {
	Model

	Name string `gorm:"name" json:"name"`
	Code string `gorm:"code" json:"code"`

	Products []Product `gorm:"foreignKey:category_id;references:id" json:"products"`
}

func (Category) TableName() string {
	return "category"
}
