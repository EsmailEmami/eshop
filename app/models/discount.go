package models

import (
	"errors"
	"time"

	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/validations"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
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

func (model DiscountReqModel) ValidateCreate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Value,
			validation.Required.Error(consts.Required),
		),
		validation.Field(&model.ProductItemID,
			validation.By(validations.ExistsInDB(&dbmodels.ProductItem{}, "id", consts.ModelProductItemNotFound)),
			validation.By(func(value interface{}) error {
				if validations.IsNil(value) {
					return nil
				}

				if !validations.IsNil(model.Code) || !validations.IsNil(model.RelatedUserID) {
					return errors.New("discount cannot have a code or user for product item")
				}

				return nil
			}),
		),
		validation.Field(&model.Code,
			validation.By(validations.Code()),
			validation.By(validations.NotExistsInDB(&dbmodels.Discount{}, "code", consts.ExistedCode)),
		),
		validation.Field(&model.ExpiresIn,
			validation.By(validations.TimeGreaterThanNow()),
		),
		validation.Field(&model.RelatedUserID,
			validation.By(validations.ExistsInDB(&dbmodels.User{}, "id", consts.ModelUserNotFound)),
			validation.By(func(value interface{}) error {
				if validations.IsNil(value) {
					return nil
				}

				if validations.IsNil(model.Code) {
					return errors.New("discount must have a code for specific user")
				}

				if !validations.IsNil(model.ProductItemID) {
					return errors.New("discount cannot have a product item for specific user")
				}

				return nil
			}),
		),
	)
}

func (model DiscountReqModel) ValidateUpdate(id uuid.UUID) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Value,
			validation.Required.Error(consts.Required),
		),
		validation.Field(&model.ProductItemID,
			validation.By(validations.ExistsInDB(&dbmodels.ProductItem{}, "id", consts.ModelProductItemNotFound)),
			validation.By(func(value interface{}) error {
				if validations.IsNil(value) {
					return nil
				}

				if !validations.IsNil(model.Code) || !validations.IsNil(model.RelatedUserID) {
					return errors.New("discount cannot have a code or user for product item")
				}

				return nil
			}),
		),
		validation.Field(&model.Code,
			validation.By(validations.Code()),
			validation.By(validations.NotExistsInDBWithID(&dbmodels.Discount{}, "code", id, consts.ExistedCode)),
		),
		validation.Field(&model.ExpiresIn,
			validation.By(validations.TimeGreaterThanNow()),
		),
		validation.Field(&model.RelatedUserID,
			validation.By(validations.ExistsInDB(&dbmodels.User{}, "id", consts.ModelUserNotFound)),
			validation.By(func(value interface{}) error {
				if validations.IsNil(value) {
					return nil
				}

				if validations.IsNil(model.Code) {
					return errors.New("discount must have a code for specific user")
				}

				if !validations.IsNil(model.ProductItemID) {
					return errors.New("discount cannot have a product item for specific user")
				}

				return nil
			}),
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

type DiscountAdminOutPutModel struct {
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

type ValidateDiscountOutPutModel struct {
	Success bool                   `json:"success"`
	ID      *uuid.UUID             `json:"id"`
	Type    *dbmodels.DiscountType `json:"type,omitempty"`
	Value   *float64               `json:"value,omitempty"`
}
