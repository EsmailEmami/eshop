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

type ProductFeatureValueReqModel struct {
	Value               string    `json:"value"`
	ProductFeatureKeyID uuid.UUID `json:"productFeatureKeyId"`
	ProductID           uuid.UUID `json:"-"`
}

func (model ProductFeatureValueReqModel) ValidateCreate(db *gorm.DB) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Value,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.ProductFeatureKeyID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {

				if !dbpkg.Exists(db, &dbmodels.ProductFeatureKey{}, "id=?", value) {
					return errors.New(consts.ModelProductFeatureKeyNotFound)
				}

				return nil
			}),
		),
		validation.Field(&model.ProductID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {

				if !dbpkg.Exists(db, &dbmodels.Product{}, "id=?", value) {
					return errors.New(consts.ModelProductNotFound)
				}

				return nil
			}),
		),
	)
}

func (model ProductFeatureValueReqModel) ValidateUpdate(db *gorm.DB) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Value,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText())),
		validation.Field(&model.ProductFeatureKeyID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {

				if !dbpkg.Exists(db, &dbmodels.ProductFeatureKey{}, "id=?", value) {
					return errors.New(consts.ModelProductFeatureKeyNotFound)
				}

				return nil
			}),
		),
		validation.Field(&model.ProductID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {

				if !dbpkg.Exists(db, &dbmodels.Product{}, "id=?", value) {
					return errors.New(consts.ModelProductNotFound)
				}

				return nil
			}),
		),
	)
}

func (model ProductFeatureValueReqModel) ToDBModel() *dbmodels.ProductFeatureValue {
	return &dbmodels.ProductFeatureValue{
		Model: dbmodels.Model{
			ID: dbmodels.NewID(),
		},
		Value:               model.Value,
		ProductFeatureKeyID: model.ProductFeatureKeyID,
		ProductID:           model.ProductID,
	}
}

func (model ProductFeatureValueReqModel) MergeWithDBData(
	dbmodel *dbmodels.ProductFeatureValue,
) {
	dbmodel.Value = model.Value
	dbmodel.ProductFeatureKeyID = model.ProductFeatureKeyID
	dbmodel.ProductID = model.ProductID
}

type ProductFeatureValueOutPutModel struct {
	ID                  *uuid.UUID `gorm:"column:id"                     json:"id"`
	CreatedAt           time.Time  `gorm:"column:created_at"             json:"createdAt"`
	UpdatedAt           time.Time  `gorm:"column:updated_at"             json:"updatedAt"`
	ProductFeatureKeyID uuid.UUID  `gorm:"column:product_feature_key_id" json:"productFeatureKeyId"`
	Value               string     `gorm:"column:value"                  json:"value"`
}
