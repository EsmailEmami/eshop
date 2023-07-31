package models

import "github.com/google/uuid"

type ProductFeatureKey struct {
	Model

	ProductFeatureCategoryID uuid.UUID               `gorm:"product_feature_category_id"                          json:"productFeatureCategoryId"`
	ProductFeatureCategory   *ProductFeatureCategory `gorm:"foreignKey:product_feature_category_id;references:id" json:"productFeatureCategory"`
	Name                     string                  `gorm:"name"                                                 json:"name"`
	Values                   []ProductFeatureValue   `gorm:"foreignKey:product_feature_key_id; references:id;"    json:"values"`
}

func (ProductFeatureKey) TableName() string {
	return "product_feature_key"
}
