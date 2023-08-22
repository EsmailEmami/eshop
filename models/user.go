package models

import (
	"github.com/esmailemami/eshop/app/consts"
	"github.com/google/uuid"
)

type User struct {
	Model
	Username   string      `gorm:"column:username"                                 json:"username"`
	Password   string      `gorm:"column:password"                                 json:"-"`
	FirstName  *string     `gorm:"column:first_name"                               json:"firstName"`
	LastName   *string     `gorm:"column:last_name"                                json:"lastName"`
	Mobile     *string     `gorm:"column:mobile"                                   json:"mobile"`
	RoleID     *uuid.UUID  `gorm:"column:role_id"                                  json:"roleId"`
	Email      *string     `gorm:"email"                                           json:"email"`
	Role       *Role       `gorm:"foreignKey:role_id; references:id"               json:"role"`
	IsSystem   bool        `gorm:"column:is_system"                                json:"isSystem"`
	Enabled    bool        `gorm:"column:enabled"                                  json:"enabled"`
	AuthTokens []AuthToken `gorm:"foreignKey:user_id;references:id"                json:"authTokens"`
	Comments   []Comment   `gorm:"foreignKey:created_by_id;references:id"          json:"comments"`
	Addresses  []Address   `gorm:"foreignKey:created_by_id;references:id"          json:"addresses"`
	Discounts  []Discount  `gorm:"foreignKey:related_user_id;references:id"        json:"discounts"`
}

func (User) TableName() string {
	return "user"
}

func (user User) Can(action string) bool {
	if user.Role == nil {
		return false
	}

	if !user.Role.Permitted(action) {
		return false
	}

	return true
}

func (user User) IsRoot() bool {
	if user.RoleID == nil {
		return false
	}

	return user.RoleID.String() == consts.ROLE_ROOT_ID
}
