package models

import (
	"errors"
	"time"

	"github.com/esmailemami/eshop/app/validations"
	"github.com/esmailemami/eshop/consts"
	dbpkg "github.com/esmailemami/eshop/db"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ColorReqModel struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	ColorHex string `json:"colorHex"`
}

func (model ColorReqModel) ValidateCreate(db *gorm.DB) error {
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
			validation.By(func(value interface{}) error {

				if dbpkg.Exists(db, &dbmodels.Color{}, "code=?", value) {
					return errors.New(consts.ExistedCode)
				}

				return nil
			}),
		),
	)
}

func (model ColorReqModel) ValidateUpdate() error {
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
		),
	)
}

func (model ColorReqModel) ToDBModel() *dbmodels.Color {
	return &dbmodels.Color{
		Model: dbmodels.Model{
			ID: dbmodels.NewID(),
		},
		Name:     model.Name,
		Code:     model.Code,
		ColorHex: model.ColorHex,
	}
}

func (model ColorReqModel) MergeWithDBData(dbmodel *dbmodels.Color) {
	dbmodel.Name = model.Name
	dbmodel.Code = model.Code
	dbmodel.ColorHex = model.ColorHex
}

type ColorOutPutModel struct {
	ID        *uuid.UUID `gorm:"id"                json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	Name      string     `gorm:"name"              json:"name"`
	Code      string     `gorm:"code"              json:"code"`
	ColorHex  string     `gorm:"color_hex"         json:"colorHex"`
}
