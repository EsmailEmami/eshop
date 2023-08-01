package models

import (
	datatypes "github.com/esmailemami/eshop/models/data_types"
	"github.com/google/uuid"
)

type Comment struct {
	Model

	Text           string                `gorm:"column:text"                           json:"text"`
	Rate           int                   `gorm:"column:rate"                           json:"rate"`
	StrengthPoints datatypes.StringArray `gorm:"column:strength_points"                json:"strengthPoints"`
	WeakPonits     datatypes.StringArray `gorm:"column:weak_ponits"                    json:"weakPonits"`
	ProductID      uuid.UUID             `gorm:"column:product_id"                     json:"productId"`
	Product        *Product              `gorm:"foreignKey:product_id; references:id;" json:"product"`
}

func (Comment) TableName() string {
	return "comment"
}
