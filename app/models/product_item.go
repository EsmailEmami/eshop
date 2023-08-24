package models

import (
	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/validations"
	"github.com/esmailemami/eshop/models"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type ProductItemReqModel struct {
	Price      float64                `json:"price"`
	Status     dbmodels.ProductStatus `json:"status"`
	ColorID    uuid.UUID              `json:"colorId"`
	ProductID  uuid.UUID              `json:"productId"`
	Quantity   int                    `json:"quantity"`
	IsMainItem bool                   `json:"isMainItem"`
}

func (model ProductItemReqModel) ValidateCreate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Price,
			validation.Required.Error(consts.Required)),
		validation.Field(&model.Quantity,
			validation.Required.Error(consts.Required)),
		validation.Field(&model.ColorID,
			validation.Required.Error(consts.Required),
			validation.By(validations.ExistsInDB(&dbmodels.Color{}, "id", consts.ModelColorNotFound)),
		),
		validation.Field(&model.ProductID,
			validation.Required.Error(consts.Required),
			validation.By(validations.ExistsInDB(&dbmodels.Product{}, "id", consts.ModelProductNotFound)),
		),
	)
}

func (model ProductItemReqModel) ValidateUpdate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Price,
			validation.Required.Error(consts.Required)),
		validation.Field(&model.Quantity,
			validation.Required.Error(consts.Required)),
		validation.Field(&model.ColorID,
			validation.Required.Error(consts.Required),
			validation.By(validations.ExistsInDB(&dbmodels.Color{}, "id", consts.ModelColorNotFound)),
		),
		validation.Field(&model.ProductID,
			validation.Required.Error(consts.Required),
			validation.By(validations.ExistsInDB(&dbmodels.Product{}, "id", consts.ModelProductNotFound)),
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

type ProductItemOutPutModel struct {
	ID           *uuid.UUID             `gorm:"column:id"                        json:"id"`
	Price        float64                `gorm:"column:price"                     json:"price"`
	Status       dbmodels.ProductStatus `gorm:"column:status"                    json:"status"`
	ColorID      uuid.UUID              `gorm:"column:color_id"                  json:"colorId"`
	ColorName    string                 `gorm:"column:color_name"                json:"color"`
	ProductID    uuid.UUID              `gorm:"column:product_id"                json:"productId"`
	ProductTitle string                 `gorm:"column:product_title"             json:"productTitle"`
	ProductCode  string                 `gorm:"column:product_code"              json:"productCode"`
	Quantity     int                    `gorm:"column:quantity"                  json:"quantity"`
}

type ProductItemInfoOutPutModel struct {
	ID                      *uuid.UUID                        `gorm:"column:id"                        json:"id"`
	Price                   float64                           `gorm:"column:price"                     json:"price"`
	Status                  dbmodels.ProductStatus            `gorm:"column:status"                    json:"status"`
	ColorID                 uuid.UUID                         `gorm:"column:color_id"                  json:"colorId"`
	ColorName               string                            `gorm:"column:color_name"                json:"color"`
	ProductID               uuid.UUID                         `gorm:"column:product_id"                json:"productId"`
	ProductTitle            string                            `gorm:"column:product_title"             json:"productTitle"`
	ProductCode             string                            `gorm:"column:product_code"              json:"productCode"`
	Quantity                int                               `gorm:"column:quantity"                  json:"quantity"`
	ProductShortDescription string                            `gorm:"column:product_short_description" json:"productShortDescription"`
	ProductDescription      string                            `gorm:"column:product_description"       json:"productDescription"`
	DiscountType            *dbmodels.DiscountType            `gorm:"column:discount_type"             json:"discountType,omitempty"`
	DiscountValue           *float64                          `gorm:"column:discount_value"            json:"discountValue"`
	DiscountQuantity        *int                              `gorm:"column:discount_quantity"         json:"discountQuantity"`
	Files                   []ProductItemFileOutPutModel      `gorm:"-"                                json:"files"`
	Features                []ProductItemCategoryFeatureModel `gorm:"-"                                json:"features"`
	Colors                  []ProductItemInfoColorOutPutModel `gorm:"-"                                json:"colors"`
}

type ProductItemCategoryFeatureModel struct {
	Category string                    `json:"category"`
	Items    []ProductItemFeatureModel `json:"items"`
}

type ProductItemFeatureModel struct {
	Key   string `gorm:"column:key"   json:"key"`
	Value string `gorm:"column:value" json:"value"`
}

type ProductItemFileOutPutModel struct {
	ID             *uuid.UUID      `gorm:"column:id"               json:"id"`
	OriginalName   string          `gorm:"column:original_name"    json:"originalName"`
	UniqueFileName string          `gorm:"column:unique_file_name" json:"uniqueFineName"`
	FileType       models.FileType `gorm:"column:file_type"        json:"fileType"`
	FileUrl        string          `gorm:"-"                       json:"fileUrl"`
}

type ProductItemInfoColorOutPutModel struct {
	Name          string    `gorm:"column:name"                   json:"name"`
	ColorHex      string    `gorm:"column:color_hex"              json:"colorHex"`
	ProductItemID uuid.UUID `gorm:"column:product_item_id"        json:"productItemId"`
}
