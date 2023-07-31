package models

import (
	"errors"
	"time"

	"github.com/esmailemami/eshop/app/validations"
	"github.com/esmailemami/eshop/consts"
	dbpkg "github.com/esmailemami/eshop/db"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductFeatureCategoryReqModel struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (model ProductFeatureCategoryReqModel) ValidateCreate(db *gorm.DB) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
			validation.By(func(value interface{}) error {

				if dbpkg.Exists(db, &dbmodels.ProductFeatureCategory{}, "name=?", value) {
					return errors.New(consts.ExistedTitle)
				}

				return nil
			}),
		),
		validation.Field(&model.Code,
			validation.Required.Error(consts.Required),
			validation.By(validations.Code()),
			validation.By(func(value interface{}) error {

				if dbpkg.Exists(db, &dbmodels.ProductFeatureCategory{}, "code=?", value) {
					return errors.New(consts.ExistedCode)
				}

				return nil
			}),
		),
	)
}

func (model ProductFeatureCategoryReqModel) ValidateUpdate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Code,
			validation.Required.Error(consts.Required),
			validation.By(validations.Code()),
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
	ID        *uuid.UUID `gorm:"id"                json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	Name      string     `gorm:"name"              json:"name"`
	Code      string     `gorm:"code"              json:"code"`
}
