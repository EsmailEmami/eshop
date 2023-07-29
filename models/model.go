package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Model مدل پایه ای جدول ها
type Model struct {
	ID          *uuid.UUID     `gorm:"primaryKey"             json:"id"`
	CreatedAt   time.Time      `gorm:"column:created_at"      json:"createdAt"`
	CreatedByID *uuid.UUID     `gorm:"column:created_by_id"   json:"createdById"`
	CreatedBy   *User          `gorm:"foreignKey:CreatedByID" json:"createdBy"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"      json:"updatedAt"`
	UpdatedBy   *User          `gorm:"foreignKey:UpdatedByID" json:"updatedBy"`
	UpdatedByID *uuid.UUID     `gorm:"column:updated_by_id"   json:"updatedById"`
	DeletedAt   gorm.DeletedAt `                              json:"-"`
	DeletedByID *uuid.UUID     `gorm:"column:deleted_by_id"   json:"-"`
	DeletedBy   *User          `gorm:"foreignKey:DeletedByID" json:"-"`
}

type BasicModel struct {
	ID        *uuid.UUID     `gorm:"primaryKey"        json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"deleted_at"        json:"-"`
}

func NewID() *uuid.UUID {
	id := uuid.New()
	return &id
}
