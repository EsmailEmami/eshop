package models

type Category struct {
	Model

	Name string `gorm:"column:name" json:"name"`
	Code string `gorm:"column:code" json:"code"`

	Products []Product `gorm:"foreignKey:category_id;references:id" json:"products"`
}

func (Category) TableName() string {
	return "category"
}
