package models

import (
	"github.com/esmailemami/eshop/consts"
	"github.com/google/uuid"
)

// User مدل اطلاعات کاربر
type User struct {
	Model
	Username   string     `gorm:"column:username" json:"username"`
	Password   string     `gorm:"column:password" json:"-"`
	FirstName  string     `gorm:"column:first_name" json:"firstName"`
	LastName   string     `gorm:"column:last_name" json:"lastName"`
	Mobile     string     `gorm:"column:mobile" json:"mobile"`
	NoeSematID *uuid.UUID `gorm:"column:noe_semat_id" json:"noeSematId"`
	NoeSemat   *NoeSemat  `gorm:"foreignKey:noe_semat_id;references:id" json:"noeSemat"`
	IsSystem   bool       `gorm:"column:is_system" json:"isSystem"`
	Enabled    bool       `gorm:"column:enabled" json:"enabled"`
}

// TableName نام جدول اطلاعات کاربر
func (User) TableName() string {
	return "user"
}

// Can check user's permission
func (user User) Can(action string) bool {
	if user.NoeSemat == nil {
		return false
	}

	if !user.NoeSemat.Permitted(action) {
		return false
	}

	return true
}

func (user User) IsRoot() bool {
	if user.NoeSematID == nil {
		return false
	}

	return user.NoeSematID.String() == consts.NOE_SEMAT_ROOT_ID
}

// UserInfo این مدل حاوی فیلدهایی اضافه برای کوئری سلکت های خاص است
type UserInfo struct {
	ID           *uuid.UUID `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"column:username" json:"username"`
	FirstName    string     `gorm:"column:first_name" json:"firstName"`
	LastName     string     `gorm:"column:last_name" json:"lastName"`
	Mobile       string     `gorm:"column:mobile" json:"mobile"`
	Enabled      bool       `gorm:"column:enabled" json:"enabled"`
	Fullname     string     `gorm:"column:full_name" json:"fullName"`
	NoeSematName string     `gorm:"column:noe_semat_name" json:"noeSematName"`
}

// TableName نام جدول اطلاعات کاربر
func (UserInfo) TableName() string {
	return "user"
}
