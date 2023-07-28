package models

import (
	"database/sql/driver"
	"encoding/json"
)

type Role struct {
	Model

	Name        string          `gorm:"column:name"        json:"name"`
	Code        string          `gorm:"column:code"        json:"code"`
	IsSystem    bool            `gorm:"column:is_system"   json:"isSystem"`
	Permissions RolePermissions `gorm:"column:permissions" json:"permissions"`
}

func (Role) TableName() string {
	return "role"
}

func (model Role) Permitted(action string) bool {
	for _, p := range model.Permissions {
		if p == action {
			return true
		}
	}

	return false
}

type RolePermissions []string

func (p RolePermissions) Value() (driver.Value, error) {
	valueString, err := json.Marshal(p)
	return string(valueString), err
}

func (j *RolePermissions) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	var bts []byte
	switch v := value.(type) {
	case []byte:
		bts = v
	case string:
		bts = []byte(v)
	case nil:
		*j = nil
		return nil
	}
	return json.Unmarshal(bts, &j)
}
