package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Discount struct {
	Model

	ProductItemID *uuid.UUID   `gorm:"column:product_item_id"                                  json:"productItemId,omitempty"`
	ProductItem   *ProductItem `gorm:"foreignKey:product_item_id;references:id"                json:"productItem"`
	Type          DiscountType `gorm:"column:type"                                             json:"type"`
	Value         float64      `form:"column:value"                                            json:"value"`
	Quantity      *int         `gorm:"column:quantity"                                         json:"quantity,omitempty"`
	ExpiresIn     *time.Time   `gorm:"column:expires_in"                                       json:"expiresIn,omitempty"`
	Code          *string      `gorm:"column:code"                                             json:"code,omitempty"`
	RelatedUserID *uuid.UUID   `gorm:"column:related_user_id"                                  json:"relatedUserId,omitempty"`
	RelatedUser   *User        `gorm:"foreignKey:related_user_id;references:id"                json:"relatedUser,omitempty"`
}

func (Discount) TableName() string {
	return "discount"
}

func (d *Discount) IsValidToUse(userID, productItemID *uuid.UUID) error {
	if d.ExpiresIn != nil && d.ExpiresIn.Before(time.Now()) {
		return errors.New("discount expired")
	}

	if d.Quantity != nil && *d.Quantity <= 0 {
		return errors.New("discount reached to full quantity limit")
	}

	if d.RelatedUserID != nil && (userID == nil || *d.RelatedUserID != *userID) {
		return errors.New("discount is not for you")
	}

	if d.ProductItemID != nil && (productItemID == nil || *d.ProductItemID != *productItemID) {
		return errors.New("discount is not for this product")
	}
	return nil
}

type DiscountType int

const (
	DiscountTypeNumeric DiscountType = iota
	DiscountTypePercent
	DiscountTypeCode
)

func (d DiscountType) String() string {
	switch d {
	case DiscountTypeNumeric:
		return "numeric"
	case DiscountTypePercent:
		return "percent"
	case DiscountTypeCode:
		return "code"
	default:
		return "unknown"
	}
}
