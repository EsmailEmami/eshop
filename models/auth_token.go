package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthToken struct {
	BasicModel
	UserID         uuid.UUID      `gorm:"column:user_id"`
	User           *User          `gorm:"foreignKey:user_id;references:id"  json:"user"`
	Revoked        bool           `gorm:"column:revoked"`
	ExpiresAt      time.Time      `gorm:"column:expires_at"`
	LoginHistories []LoginHistory `gorm:"foreignKey:token_id;references:id" json:"loginHistories"`
}

func (model AuthToken) TableName() string {
	return "auth_token"
}
