package models

import (
	"errors"
	"time"

	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/validations"
	dbpkg "github.com/esmailemami/eshop/db"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductFeatureKeyReqModel struct {
	Name                     string    `json:"name"`
	ProductFeatureCategoryID uuid.UUID `json:"productFeatureCategoryId"`
}

func (model ProductFeatureKeyReqModel) ValidateCreate(db *gorm.DB) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
			validation.By(func(value interface{}) error {

				if dbpkg.Exists(db, &dbmodels.ProductFeatureKey{}, "name=?", value) {
					return errors.New(consts.ExistedTitle)
				}

				return nil
			}),
		),
		validation.Field(&model.ProductFeatureCategoryID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {

				if !dbpkg.Exists(db, &dbmodels.ProductFeatureCategory{}, "id=?", value) {
					return errors.New(consts.ModelCategoryNotFound)
				}

				return nil
			}),
		),
	)
}

func (model ProductFeatureKeyReqModel) ValidateUpdate(db *gorm.DB) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText())),
		validation.Field(&model.ProductFeatureCategoryID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {

				if !dbpkg.Exists(db, &dbmodels.ProductFeatureCategory{}, "id=?", value) {
					return errors.New(consts.ModelCategoryNotFound)
				}

				return nil
			}),
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
