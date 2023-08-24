package models

import (
	"time"

	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/validations"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type ProductFeatureKeyReqModel struct {
	Name                     string    `json:"name"`
	ProductFeatureCategoryID uuid.UUID `json:"productFeatureCategoryId"`
}

func (model ProductFeatureKeyReqModel) ValidateCreate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
			validation.By(validations.NotExistsInDB(&dbmodels.ProductFeatureKey{}, "name", consts.ExistedTitle)),
		),
		validation.Field(&model.ProductFeatureCategoryID,
			validation.Required.Error(consts.Required),
			validation.By(validations.ExistsInDB(&dbmodels.ProductFeatureCategory{}, "id", consts.ModelProductFeatureCategoryNotFound)),
		),
	)
}

func (model ProductFeatureKeyReqModel) ValidateUpdate(id uuid.UUID) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
			validation.By(validations.NotExistsInDBWithID(&dbmodels.ProductFeatureKey{}, "name", id, consts.ExistedTitle)),
		),
		validation.Field(&model.ProductFeatureCategoryID,
			validation.Required.Error(consts.Required),
			validation.By(validations.ExistsInDB(&dbmodels.ProductFeatureCategory{}, "id", consts.ModelProductFeatureCategoryNotFound)),
		),
	)
}

func (model ProductFeatureKeyReqModel) ToDBModel() *dbmodels.ProductFeatureKey {
	return &dbmodels.ProductFeatureKey{
		Model: dbmodels.Model{
			ID: dbmodels.NewID(),
		},
		Name:                     model.Name,
		ProductFeatureCategoryID: model.ProductFeatureCategoryID,
	}
}

func (model ProductFeatureKeyReqModel) MergeWithDBData(
	dbmodel *dbmodels.ProductFeatureKey,
) {
	dbmodel.Name = model.Name
	dbmodel.ProductFeatureCategoryID = model.ProductFeatureCategoryID
}

type ProductFeatureKeyOutPutModel struct {
	ID        *uuid.UUID `gorm:"column:id"         json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	Name      string     `gorm:"column:name"       json:"name"`
}
