package models

import (
	"database/sql/driver"
	"encoding/json"
)

// NoeSemat مدل اطلاعات، انواع سمت کاربران
type NoeSemat struct {
	Model

	Name        string              `gorm:"column:name"        json:"name"`
	Code        string              `gorm:"column:code"        json:"code"`
	IsSystem    bool                `gorm:"column:is_system"   json:"isSystem"`
	Permissions NoeSematPermissions `gorm:"column:permissions" json:"permissions"`
}

// TableName نام جدول، انواع سمت کاربران
func (NoeSemat) TableName() string {
	return "noe_semat"
}

func (model NoeSemat) Permitted(action string) bool {
	for _, p := range model.Permissions {
		if p == action {
			return true
		}
	}

	return false
}

type NoeSematPermissions []string

func (p NoeSematPermissions) Value() (driver.Value, error) {
	valueString, err := json.Marshal(p)
	return string(valueString), err
}

func (j *NoeSematPermissions) Scan(value interface{}) error {
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
