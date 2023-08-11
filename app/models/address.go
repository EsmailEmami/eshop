package models

import (
	"strings"
	"time"

	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/validations"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AddressReqModel struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Plaque       int    `json:"plaque"`
	PhoneNumber  string `json:"phoneNumber"`
	NationalCode string `json:"nationalCode"`
	Address      string `json:"address"`
	PostalCode   string `json:"postalCode"`
}

func (model AddressReqModel) ValidateCreate(db *gorm.DB) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.FirstName,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.LastName,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.PhoneNumber,
			validation.Required.Error(consts.Required),
			validation.By(validations.IsValidMobileNumber()),
		),
		validation.Field(&model.NationalCode,
			validation.By(func(value interface{}) error {
				if value != nil && strings.TrimSpace(value.(string)) != "" {
					fn := validations.IsValidNationalCode()
					return fn(value)
				}

				return nil
			}),
		),
		validation.Field(&model.LastName,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.PostalCode,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
			validation.By(validations.IsValidPostalCode()),
		),
	)
}

func (model AddressReqModel) ValidateUpdate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.FirstName,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.LastName,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.PhoneNumber,
			validation.Required.Error(consts.Required),
			validation.By(validations.IsValidMobileNumber()),
		),
		validation.Field(&model.NationalCode,
			validation.By(func(value interface{}) error {
				if value != nil && strings.TrimSpace(value.(string)) != "" {
					fn := validations.IsValidNationalCode()
					return fn(value)
				}

				return nil
			}),
		),
		validation.Field(&model.LastName,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.PostalCode,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
			validation.By(validations.IsValidPostalCode()),
		),
	)
}

func (model *AddressReqModel) ToDBModel() *dbmodels.Address {
	return &dbmodels.Address{
		Model: dbmodels.Model{
			ID: dbmodels.NewID(),
		},
		FirstName:    model.FirstName,
		LastName:     model.LastName,
		Plaque:       model.Plaque,
		PhoneNumber:  model.PhoneNumber,
		NationalCode: model.NationalCode,
		Address:      model.Address,
		PostalCode:   model.PostalCode,
	}
}

func (model AddressReqModel) MergeWithDBData(dbmodel *dbmodels.Address) {
	dbmodel.FirstName = model.FirstName
	dbmodel.LastName = model.LastName
	dbmodel.Plaque = model.Plaque
	dbmodel.PhoneNumber = model.PhoneNumber
	dbmodel.NationalCode = model.NationalCode
	dbmodel.Address = model.Address
	dbmodel.PostalCode = model.PostalCode
}

type AddressOutPutModel struct {
	ID           *uuid.UUID `gorm:"column:id"         json:"id"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	FirstName    string     `gorm:"column:first_name" json:"firstName"`
	LastName     string     `gorm:"column:last_name"         json:"lastName"`
	Plaque       int        `gorm:"column:plaque"            json:"plaque"`
	PhoneNumber  string     `gorm:"column:phone_number"      json:"phoneNumber"`
	NationalCode string     `gorm:"column:national_code"     json:"nationalCode"`
	PostalCode   string     `gorm:"column:postal_code"     json:"postalCode"`
	Address      string     `gorm:"column:address"              json:"address"`
}
