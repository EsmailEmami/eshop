package models

import "github.com/google/uuid"

type ProductFeatureValue struct {
	Model

	ProductFeatureKeyID uuid.UUID          `gorm:"column:product_feature_key_id"                     json:"productFeatureKeyId"`
	ProductFeatureKey   *ProductFeatureKey `gorm:"foreignKey:product_feature_key_id; references:id;" json:"productFeatureKey"`
	ProductID           uuid.UUID          `gorm:"column:product_id"                                 json:"productId"`
	Product             *Product           `gorm:"foreignKey:product_id; references:id;"             json:"product"`
	Value               string             `gorm:"column:value"                                      json:"value"`
}

func (ProductFeatureValue) TableName() string {
	return "product_feature_value"
}
