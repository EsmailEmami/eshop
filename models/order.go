package models

import "time"

type Order struct {
	Model

	Code          string        `gorm:"column:code"                       json:"code"`
	Items         []OrderItem   `gorm:"foreignKey:order_id;references:id" json:"items"`
	Status        OrderStatus   `gorm:"column:status"                     json:"status"`
	Price         float64       `gorm:"column:price"                      json:"price"`
	PaidAt        *time.Time    `gorm:"column:paid_at"                    json:"paidAt"`
	TotalPrice    float64       `gorm:"column:total_price"                json:"totalPrice"`
	DiscountPrice *float64      `gorm:"column:discount_price"             json:"discountPrice,omitempty"`
	DiscountValue *float64      `gorm:"column:discount_value"             json:"discountValue,omitempty"`
	DiscountType  *DiscountType `gorm:"column:discount_type"              json:"discountType,omitempty"`

	// keep the address
	FirstName    string `gorm:"column:first_name"        json:"firstName"`
	LastName     string `gorm:"column:last_name"         json:"lastName"`
	Plaque       int    `gorm:"column:plaque"            json:"plaque"`
	PhoneNumber  string `gorm:"column:phone_number"      json:"phoneNumber"`
	NationalCode string `gorm:"column:national_code"     json:"nationalCode"`
	PostalCode   string `gorm:"column:postal_code"       json:"postalCode"`
	Address      string `gorm:"column:address"           json:"address"`
}

func (Order) TableName() string {
	return "order"
}

type OrderStatus int

const (
	OrderStatusOpen = iota
	OrderStatusPaid
	OrderStatusProcessing
	OrderStatusSent
	OrderStatusReceived
)

func (os OrderStatus) String() string {
	switch os {
	case OrderStatusOpen:
		return "Open"
	case OrderStatusPaid:
		return "Paid"
	case OrderStatusProcessing:
		return "Processing"
	case OrderStatusSent:
		return "Sent"
	case OrderStatusReceived:
		return "received"
	default:
		return "unknown"
	}
}
