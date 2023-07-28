package models

type Example struct {
	Model

	Name string `gorm:"column:name" json:"name,omitempty"`
	Code string `gorm:"column:code" json:"code,omitempty"`
}

func (Example) TableName() string {
	return "--"
}
