package models

type Order struct {
	Model

	Code   string      `gorm:"code"                              json:"code"`
	Items  []OrderItem `gorm:"foreignKey:order_id;references:id" json:"items"`
	Status OrderStatus `gorm:"status"                            json:"status"`
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
