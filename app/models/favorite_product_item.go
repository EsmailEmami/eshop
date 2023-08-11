package models

import (
	"errors"

	"github.com/esmailemami/eshop/app/consts"
	dbpkg "github.com/esmailemami/eshop/db"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FavoriteProductItemReqModel struct {
	ProductItemID uuid.UUID `json:"productItemId"`
}

func (model FavoriteProductItemReqModel) ValidateCreate(db *gorm.DB) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.ProductItemID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {

				if !dbpkg.Exists(db, &dbmodels.ProductItem{}, "id=?", value) {
					return errors.New(consts.ModelProductNotFound)
				}

				return nil
			}),
		),
	)
}

func (model *FavoriteProductItemReqModel) ToDBModel() *dbmodels.FavoriteProductItem {
	return &dbmodels.FavoriteProductItem{
		Model: dbmodels.Model{
			ID: dbmodels.NewID(),
		},
		ProductItemID: model.ProductItemID,
	}
}
