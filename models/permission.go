package models

type Permissions struct {
	Model

	GroupName string `gorm:"column:group_name" json:"groupName,omitempty"`
	Name      string `gorm:"column:name"       json:"name,omitempty"`
	Code      string `gorm:"column:code"       json:"code,omitempty"`
}

func (Permissions) TableName() string {
	return "permissions"
}
