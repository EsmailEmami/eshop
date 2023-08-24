package models

import (
	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/validations"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type FavoriteProductItemReqModel struct {
	ProductItemID uuid.UUID `json:"productItemId"`
}

func (model FavoriteProductItemReqModel) ValidateCreate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.ProductItemID,
			validation.Required.Error(consts.Required),
			validation.By(validations.ExistsInDB(&dbmodels.ProductItem{}, "id", consts.ModelProductItemNotFound)),
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
