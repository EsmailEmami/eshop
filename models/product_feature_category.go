package models

type ProductFeatureCategory struct {
	Model

	Name string              `gorm:"column:name"                                            json:"name"`
	Code string              `gorm:"column:code"                                            json:"code"`
	Keys []ProductFeatureKey `gorm:"foreignKey:product_feature_category_id; references:id;" json:"keys"`
}

func (ProductFeatureCategory) TableName() string {
	return "product_feature_category"
}
