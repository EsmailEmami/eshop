package models

import (
	"time"

	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/validations"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type ProductFeatureCategoryReqModel struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (model ProductFeatureCategoryReqModel) ValidateCreate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
			validation.By(validations.NotExistsInDB(&dbmodels.ProductFeatureCategory{}, "name", consts.ExistedTitle)),
		),
		validation.Field(&model.Code,
			validation.Required.Error(consts.Required),
			validation.By(validations.Code()),
			validation.By(validations.NotExistsInDB(&dbmodels.ProductFeatureCategory{}, "code", consts.ExistedCode)),
		),
	)
}

func (model ProductFeatureCategoryReqModel) ValidateUpdate(id uuid.UUID) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
			validation.By(validations.NotExistsInDBWithID(&dbmodels.ProductFeatureCategory{}, "name", id, consts.ExistedTitle)),
		),
		validation.Field(&model.Code,
			validation.Required.Error(consts.Required),
			validation.By(validations.Code()),
			validation.By(validations.NotExistsInDBWithID(&dbmodels.ProductFeatureCategory{}, "code", id, consts.ExistedCode)),
		),
	)
}

func (model ProductFeatureCategoryReqModel) ToDBModel() *dbmodels.ProductFeatureCategory {
	return &dbmodels.ProductFeatureCategory{
		Model: dbmodels.Model{
			ID: dbmodels.NewID(),
		},
		Name: model.Name,
		Code: model.Code,
	}
}

func (model ProductFeatureCategoryReqModel) MergeWithDBData(
	dbmodel *dbmodels.ProductFeatureCategory,
) {
	dbmodel.Name = model.Name
	dbmodel.Code = model.Code
}

type ProductFeatureCategoryOutPutModel struct {
	ID        *uuid.UUID `gorm:"column:id"         json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	Name      string     `gorm:"column:name"       json:"name"`
	Code      string     `gorm:"column:code"       json:"code"`
}
