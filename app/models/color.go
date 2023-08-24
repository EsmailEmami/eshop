package models

import (
	"time"

	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/validations"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type ColorReqModel struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	ColorHex string `json:"colorHex"`
}

func (model ColorReqModel) ValidateCreate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.ColorHex,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Code,
			validation.Required.Error(consts.Required),
			validation.By(validations.Code()),
			validation.By(validations.NotExistsInDB(&dbmodels.Color{}, "code", consts.ExistedCode)),
		),
	)
}

func (model ColorReqModel) ValidateUpdate(id uuid.UUID) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.ColorHex,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Code,
			validation.Required.Error(consts.Required),
			validation.By(validations.Code()),
			validation.By(validations.NotExistsInDBWithID(&dbmodels.Color{}, "code", id, consts.ExistedCode)),
		),
	)
}

func (model *ColorReqModel) ToDBModel() *dbmodels.Color {
	return &dbmodels.Color{
		Model: dbmodels.Model{
			ID: dbmodels.NewID(),
		},
		Name:     model.Name,
		Code:     model.Code,
		ColorHex: model.ColorHex,
	}
}

func (model *ColorReqModel) MergeWithDBData(dbmodel *dbmodels.Color) {
	dbmodel.Name = model.Name
	dbmodel.Code = model.Code
	dbmodel.ColorHex = model.ColorHex
}

type ColorOutPutModel struct {
	ID        *uuid.UUID `gorm:"column:id"         json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	Name      string     `gorm:"column:name"       json:"name"`
	Code      string     `gorm:"column:code"       json:"code"`
	ColorHex  string     `gorm:"column:color_hex"  json:"colorHex"`
}
