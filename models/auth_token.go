package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthToken struct {
	BasicModel
	UserID         uuid.UUID      `gorm:"column:user_id"                    json:"userId"`
	User           *User          `gorm:"foreignKey:user_id;references:id"  json:"user"`
	Revoked        bool           `gorm:"column:revoked"                    json:"revoked"`
	ExpiresAt      time.Time      `gorm:"column:expires_at"                 json:"expiresAt"`
	LoginHistories []LoginHistory `gorm:"foreignKey:token_id;references:id" json:"loginHistories"`
}

func (model AuthToken) TableName() string {
	return "auth_token"
}
