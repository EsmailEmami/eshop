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
	ProductItemID *uuid.UUID            `json:"productItemId,omitempty"`
	Type          dbmodels.DiscountType `json:"type"`
	Value         float64               `json:"value"`
	Quantity      *int                  `json:"quantity,omitempty"`
	ExpiresIn     *time.Time            `json:"expiresIn,omitempty"`
	Code          *string               `json:"code,omitempty"`
	RelatedUserID *uuid.UUID            `json:"relatedUserId,omitempty"`
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
	ID                  *uuid.UUID            `gorm:"column:id"                                               json:"id"`
	CreatedAt           time.Time             `gorm:"column:created_at"                                       json:"createdAt"`
	UpdatedAt           time.Time             `gorm:"column:updated_at"                                       json:"updatedAt"`
	ProductItemID       *uuid.UUID            `gorm:"column:product_item_id"                                  json:"productItemId,omitempty"`
	ProductName         *string               `gorm:"column:product_name"                                     json:"productName,omitempty"`
	Type                dbmodels.DiscountType `gorm:"column:type"                                             json:"type"`
	TypeName            string                `gorm:"column:type_name"                                        json:"typeName"`
	Value               float64               `form:"column:value"                                            json:"value"`
	Quantity            *int                  `gorm:"column:quantity"                                         json:"quantity,omitempty"`
	ExpiresIn           *time.Time            `gorm:"column:expires_in"                                       json:"expiresIn,omitempty"`
	Code                *string               `gorm:"column:code"                                             json:"code,omitempty"`
	CreatorUserID       *uuid.UUID            `gorm:"column:creator_user_id"                                  json:"creatorUserId,omitempty"`
	CreatorUsername     *string               `gorm:"column:creator_user_username"                            json:"creatorUsername,omitempty"`
	RelatedUserID       *uuid.UUID            `gorm:"column:related_user_id"                                  json:"relatedUserId,omitempty"`
	RelatedUserUsername *string               `gorm:"column:related_user_username"                            json:"relatedUserUsername,omitempty"`
}
