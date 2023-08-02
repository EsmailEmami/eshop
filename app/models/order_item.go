package models

import (
	"errors"

	"github.com/esmailemami/eshop/consts"
	dbpkg "github.com/esmailemami/eshop/db"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderItemReqModel struct {
	OrderID       uuid.UUID `json:"-"`
	ProductItemID uuid.UUID `json:"productItemId"`
	Quantity      int       `json:"quantity"`
	Price         float64   `json:"-"`
}

func (model OrderItemReqModel) ValidateCreate(db *gorm.DB) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Quantity,
			validation.Required.Error(consts.Required),
		),
		validation.Field(&model.ProductItemID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {

				if !dbpkg.Exists(db, &dbmodels.ProductItem{}, "id=?", value) {
					return errors.New("کالای مورد نظر یافت نشد")
				}

				return nil
			}),
		),
	)
}

func (model OrderItemReqModel) ValidateUpdate(db *gorm.DB) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Quantity,
			validation.Required.Error(consts.Required),
		),
		validation.Field(&model.ProductItemID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {

				if !dbpkg.Exists(db, &dbmodels.ProductItem{}, "id=?", value) {
					return errors.New("کالای مورد نظر یافت نشد")
				}

				return nil
			}),
		),
	)
}

func (model *OrderItemReqModel) ToDBModel() *dbmodels.OrderItem {
	return &dbmodels.OrderItem{
		Model: dbmodels.Model{
			ID: dbmodels.NewID(),
		},
		ProductItemID: model.ProductItemID,
		OrderID:       model.OrderID,
		Quantity:      model.Quantity,
		Price:         model.Price,
	}
}

func (model OrderItemReqModel) MergeWithDBData(dbmodel *dbmodels.OrderItem) {
	dbmodel.ProductItemID = model.ProductItemID
	dbmodel.OrderID = model.OrderID
	dbmodel.Quantity = model.Quantity
	dbmodel.Price = model.Price
}
