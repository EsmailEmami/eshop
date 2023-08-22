package models

import (
	"time"

	"github.com/esmailemami/eshop/app/consts"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscountReqModel struct {
	ProductItemID *uuid.UUID            `gorm:"column:product_item_id"                                  json:"productItemId,omitempty"`
	Type          dbmodels.DiscountType `gorm:"column:type"                                             json:"type"`
	Value         float64               `form:"column:value"                                            json:"value"`
	Quantity      *int                  `gorm:"column:quantity"                                         json:"quantity,omitempty"`
	ExpiresIn     *time.Time            `gorm:"column:expires_in"                                       json:"expiresIn,omitempty"`
	Code          *string               `gorm:"column:code"                                             json:"code,omitempty"`
	RelatedUserID *uuid.UUID            `gorm:"column:related_user_id"                                  json:"relatedUserId,omitempty"`
}

func (model DiscountReqModel) ValidateCreate(db *gorm.DB) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Value,
			validation.Required.Error(consts.Required),
		),
	)
}

func (model DiscountReqModel) ValidateUpdate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Value,
			validation.Required.Error(consts.Required),
		),
	)
}

func (model *DiscountReqModel) ToDBModel() *dbmodels.Discount {
	return &dbmodels.Discount{
		Model: dbmodels.Model{
			ID: dbmodels.NewID(),
		},
		ProductItemID: model.ProductItemID,
		Type:          model.Type,
		Value:         model.Value,
		Quantity:      model.Quantity,
		ExpiresIn:     model.ExpiresIn,
		Code:          model.Code,
		RelatedUserID: model.RelatedUserID,
	}
}

func (model *DiscountReqModel) MergeWithDBData(dbmodel *dbmodels.Discount) {
	dbmodel.ProductItemID = model.ProductItemID
	dbmodel.Type = model.Type
	dbmodel.Value = model.Value
	dbmodel.Quantity = model.Quantity
	dbmodel.ExpiresIn = model.ExpiresIn
	dbmodel.Code = model.Code
	dbmodel.RelatedUserID = model.RelatedUserID
}

type DiscountOutPutModel struct {
	ID          *uuid.UUID `gorm:"column:id"         json:"id"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	Name        string     `gorm:"column:name"       json:"name"`
	Code        string     `gorm:"column:code"       json:"code"`
	DiscountHex string     `gorm:"column:Discount_hex"  json:"DiscountHex"`
}
