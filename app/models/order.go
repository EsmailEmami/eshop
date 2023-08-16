package models

import (
	"time"

	dbmodels "github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
)

type OrderOutPutModel struct {
	Price float64                `gorm:"-" json:"price"`
	Items []OrderItemOutPutModel `gorm:"-" json:"items"`
}

type AdminOrderOutPutModel struct {
	UserID    uuid.UUID            `gorm:"column:user_id"         json:"userId"`
	OrderID   uuid.UUID            `gorm:"column:order_id"        json:"orderId"`
	Username  string               `gorm:"column:username"        json:"username"`
	Status    dbmodels.OrderStatus `gorm:"column:status"          json:"status"`
	Price     float64              `gorm:"column:price"           json:"price"`
	CreatedAt time.Time            `gorm:"column:created_at"      json:"createdAt"`
	UpdatedAt time.Time            `gorm:"column:updated_at"      json:"updatedAt"`
}
