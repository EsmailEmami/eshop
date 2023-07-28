package models

import (
	"github.com/google/uuid"
)

// LoginHistory
type LoginHistory struct {
	BasicModel
	IP        *string    `gorm:"column:ip" json:"ip"`
	UserAgent *string    `gorm:"column:user_agent" json:"userAgent"`
	UserID    uuid.UUID  `gorm:"column:user_id" json:"userId"`
	User      *User      `gorm:"foreignKey:user_id;references:id" json:"user"`
	TokenID   *uuid.UUID `gorm:"column:token_id" json:"tokenId"`
	Token     *AuthToken `gorm:"foreignKey:token_id;references:id" json:"token"`
}

func (LoginHistory) TableName() string {
	return "login_history"
}
