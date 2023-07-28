package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthToken struct {
	ID        uuid.UUID `gorm:"column:id"`
	UserID    uuid.UUID `gorm:"column:user_id"`
	Revoked   bool      `gorm:"column:revoked"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
}

func (model AuthToken) TableName() string {
	return "auth_token"
}
