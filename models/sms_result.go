package models

import (
	"time"

	"github.com/esmailemami/eshop/app/services/notifier/sms"
	"github.com/google/uuid"
)

type SmsResult struct {
	ID        *uuid.UUID `gorm:"column:primaryKey"        json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"createdAt"`

	sms.SmsResult
}

func (SmsResult) TableName() string {
	return "sms_result"
}
