package models

import (
	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/validations"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type OrderItemReqModel struct {
	OrderID       uuid.UUID `json:"-"`
	ProductItemID uuid.UUID `json:"productItemId"`
	Quantity      int       `json:"quantity"`
	Price         float64   `json:"-"`
}

func (model OrderItemReqModel) ValidateCreate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Quantity,
			validation.Required.Error(consts.Required),
		),
		validation.Field(&model.ProductItemID,
			validation.Required.Error(consts.Required),
			validation.By(validations.ExistsInDB(&dbmodels.ProductItem{}, "id", consts.ModelProductItemNotFound)),
		),
	)
}

func (model OrderItemReqModel) ValidateUpdate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Quantity,
			validation.Required.Error(consts.Required),
		),
		validation.Field(&model.ProductItemID,
			validation.Required.Error(consts.Required),
			validation.By(validations.ExistsInDB(&dbmodels.ProductItem{}, "id", consts.ModelProductItemNotFound)),
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

type OrderItemOutPutModel struct {
	ID               uuid.UUID              `gorm:"column:id"                      json:"id"`
	ProductItemID    uuid.UUID              `gorm:"column:product_item_id"         json:"productItemId"`
	ProductName      string                 `gorm:"column:product_name"            json:"productName"`
	Price            float64                `gorm:"column:price"                   json:"price"`
	Quantity         int                    `gorm:"column:quantity"                json:"quantity"`
	FileType         dbmodels.FileType      `gorm:"column:file_type"               json:"-"`
	FileName         string                 `gorm:"column:file_name"               json:"-"`
	FileUrl          string                 `gorm:"column:file_url"                json:"fileUrl"`
	TotalPrice       float64                `gorm:"-"                              json:"totalPrice"`
	DiscountType     *dbmodels.DiscountType `gorm:"column:discount_type"           json:"discountType,omitempty"`
	DiscountValue    *float64               `gorm:"column:discount_value"          json:"discountValue"`
	DiscountQuantity *int                   `gorm:"column:discount_quantity"       json:"discountQuantity"`
}
