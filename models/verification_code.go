package models

import "time"

type VerificationCode struct {
	BasicModel

	ExpireAt   time.Time             `gorm:"column:expire_at" json:"expireAt"`
	MaxRetires int                   `gorm:"column:max_retires" json:"maxRetires"`
	Scope      VerificationCodeScope `gorm:"column:scope" json:"scope"`
	Key        string                `gorm:"column:key" json:"key"`
	Value      string                `gorm:"column:value" json:"value"`
	Verified   bool                  `gorm:"column:verified" json:"verified"`
	Attempts   int                   `gorm:"column:attempts" json:"attempts"`
}

func (VerificationCode) TableName() string {
	return "verification_code"
}

type VerificationCodeScope int

const (
	VerificationCodeScopeEmail VerificationCodeScope = iota
	VerificationCodeScopeMobile
)
