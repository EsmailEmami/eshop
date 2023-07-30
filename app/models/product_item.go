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

type ProductItemReqModel struct {
	Price      float64                `json:"price"`
	Status     dbmodels.ProductStatus `json:"status"`
	ColorID    uuid.UUID              `json:"colorId"`
	ProductID  uuid.UUID              `json:"productId"`
	Quantity   int                    `json:"quantity"`
	IsMainItem bool                   `json:"isMainItem"`
}

func (model ProductItemReqModel) ValidateCreate(db *gorm.DB) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Price,
			validation.Required.Error(consts.Required)),
		validation.Field(&model.Quantity,
			validation.Required.Error(consts.Required)),
		validation.Field(&model.ColorID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {
				if !dbpkg.Exists(db, &dbmodels.Color{}, "id=?", value) {
					return errors.New("رنگ مورد نظر یافت نشد.")
				}

				return nil
			}),
		),
		validation.Field(&model.ProductID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {
				if !dbpkg.Exists(db, &dbmodels.Product{}, "id=?", value) {
					return errors.New("کالا مورد نظر یافت نشد.")
				}

				return nil
			}),
		),
	)
}

func (model ProductItemReqModel) ValidateUpdate(db *gorm.DB) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Price,
			validation.Required.Error(consts.Required)),
		validation.Field(&model.Quantity,
			validation.Required.Error(consts.Required)),
		validation.Field(&model.ColorID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {
				if !dbpkg.Exists(db, &dbmodels.Color{}, "id=?", value) {
					return errors.New("رنگ مورد نظر یافت نشد.")
				}

				return nil
			}),
		),
		validation.Field(&model.ProductID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {
				if !dbpkg.Exists(db, &dbmodels.Product{}, "id=?", value) {
					return errors.New("کالا مورد نظر یافت نشد.")
				}

				return nil
			}),
		),
	)
}

func (model *ProductItemReqModel) ToDBModel() *dbmodels.ProductItem {
	return &dbmodels.ProductItem{
		Model: dbmodels.Model{
			ID: dbmodels.NewID(),
		},
		Price:          model.Price,
		Status:         model.Status,
		ColorID:        model.ColorID,
		ProductID:      model.ProductID,
		BoughtQuantity: 0,
		Quantity:       model.Quantity,
	}
}

func (model ProductItemReqModel) MergeWithDBData(dbmodel *dbmodels.ProductItem) {
	dbmodel.Price = model.Price
	dbmodel.Status = model.Status
	dbmodel.ColorID = model.ColorID
	dbmodel.ProductID = model.ProductID
	dbmodel.Quantity = model.Quantity
}

type ProductItemInfoOutPutModel struct {
	ID                      *uuid.UUID             `gorm:"id"                        json:"id"`
	Price                   float64                `gorm:"price"                     json:"price"`
	Status                  dbmodels.ProductStatus `gorm:"status"                    json:"status"`
	ColorID                 uuid.UUID              `gorm:"color_id"                  json:"colorId"`
	ColorName               string                 `gorm:"color_name"                json:"color"`
	ProductID               uuid.UUID              `gorm:"product_id"                json:"productId"`
	ProductTitle            string                 `gorm:"product_title"             json:"productTitle"`
	ProductCode             string                 `gorm:"product_code"              json:"productCode"`
	Quantity                int                    `gorm:"quantity"                  json:"quantity"`
	ProductShortDescription string                 `gorm:"product_short_description" json:"productShortDescription"`
	ProductDescription      string                 `gorm:"product_description"       json:"productDescription"`
}
