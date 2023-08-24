package models

import (
	"time"

	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/validations"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type CategoryReqModel struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (model CategoryReqModel) ValidateCreate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Code,
			validation.Required.Error(consts.Required),
			validation.By(validations.Code()),
			validation.By(validations.NotExistsInDB(&dbmodels.Category{}, "code", consts.ExistedCode)),
		),
	)
}

func (model CategoryReqModel) ValidateUpdate(id uuid.UUID) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Code,
			validation.Required.Error(consts.Required),
			validation.By(validations.Code()),
			validation.By(validations.NotExistsInDBWithID(&dbmodels.Category{}, "code", id, consts.ExistedCode)),
		),
	)
}

func (model CategoryReqModel) ToDBModel() *dbmodels.Category {
	return &dbmodels.Category{
		Model: dbmodels.Model{
			ID: dbmodels.NewID(),
		},
		Name: model.Name,
		Code: model.Code,
	}
}

func (model CategoryReqModel) MergeWithDBData(dbmodel *dbmodels.Category) {
	dbmodel.Name = model.Name
	dbmodel.Code = model.Code
}

type CategoryOutPutModel struct {
	ID        *uuid.UUID `gorm:"column:id"         json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	Name      string     `gorm:"column:name"       json:"name"`
	Code      string     `gorm:"column:code"       json:"code"`
}
