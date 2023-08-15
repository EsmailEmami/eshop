package models

import (
	"time"

	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
)

type UserDashboardInfoOutPutModel struct {
	ID        *uuid.UUID `gorm:"column:id"                              json:"id"`
	Username  string     `gorm:"column:username"                        json:"username"`
	FirstName *string    `gorm:"column:first_name"                      json:"firstName"`
	LastName  *string    `gorm:"column:last_name"                       json:"lastName"`
	Mobile    *string    `gorm:"column:mobile"                          json:"mobile"`
	RoleName  string     `gorm:"column:role_name"                       json:"roleName"`
}

type UserOrderOutPutModel struct {
	ID        *uuid.UUID         `gorm:"column:id"                         json:"id"`
	CreatedAt time.Time          `gorm:"column:created_at"                 json:"createdAt"`
	Status    models.OrderStatus `gorm:"column:status"                     json:"status"`
	Price     float64            `gorm:"column:price"                      json:"price"`
	Code      string             `gorm:"column:code"                       json:"code"`
	FileUrls  []string           `gorm:"-"                                 json:"fileUrls"`
}
