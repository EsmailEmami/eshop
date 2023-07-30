package models

import (
	"errors"

	"github.com/esmailemami/eshop/app/validations"
	"github.com/esmailemami/eshop/consts"
	dbpkg "github.com/esmailemami/eshop/db"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductReqModel struct {
	Name             string    `json:"name"`
	Code             string    `json:"code"`
	BrandID          uuid.UUID `json:"brandId"`
	CategoryID       uuid.UUID `json:"categoryId"`
	Description      string    `json:"description"`
	ShortDescription string    `json:"shortDescription"`
}

func (model ProductReqModel) ValidateCreate(db *gorm.DB) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.ShortDescription,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Description,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Code,
			validation.Required.Error(consts.Required),
			validation.By(validations.Code()),
			validation.By(func(value interface{}) error {

				if dbpkg.Exists(db, &dbmodels.Product{}, "code=?", value) {
					return errors.New(consts.ExistedCode)
				}

				return nil
			}),
		),
		validation.Field(&model.BrandID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {

				if !dbpkg.Exists(db, &dbmodels.Brand{}, "id=?", value) {
					return errors.New("برند مورد نظر یافت نشد")
				}

				return nil
			}),
		),
		validation.Field(&model.CategoryID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {

				if !dbpkg.Exists(db, &dbmodels.Category{}, "id=?", value) {
					return errors.New("دسته بندی مورد نظر یافت نشد")
				}

				return nil
			}),
		),
	)
}

func (model ProductReqModel) ValidateUpdate(db *gorm.DB) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.ShortDescription,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Description,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Code,
			validation.Required.Error(consts.Required),
			validation.By(validations.Code()),
			validation.By(func(value interface{}) error {

				if dbpkg.Exists(db, &dbmodels.Product{}, "code=?", value) {
					return errors.New(consts.ExistedCode)
				}

				return nil
			}),
		),
		validation.Field(&model.BrandID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {

				if !dbpkg.Exists(db, &dbmodels.Brand{}, "id=?", value) {
					return errors.New("برند مورد نظر یافت نشد")
				}

				return nil
			}),
		),
		validation.Field(&model.CategoryID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {

				if !dbpkg.Exists(db, &dbmodels.Category{}, "id=?", value) {
					return errors.New("دسته بندی مورد نظر یافت نشد")
				}

				return nil
			}),
		),
	)
}

func (model *ProductReqModel) ToDBModel() *dbmodels.Product {
	return &dbmodels.Product{
		Model: dbmodels.Model{
			ID: dbmodels.NewID(),
		},
		Name:             model.Name,
		Code:             model.Code,
		BrandID:          model.BrandID,
		CategoryID:       model.CategoryID,
		ShortDescription: model.ShortDescription,
		Description:      model.Description,
	}
}

func (model *ProductReqModel) MergeWithDBData(dbmodel *dbmodels.Product) {
	dbmodel.Name = model.Name
	dbmodel.Code = model.Code
	dbmodel.BrandID = model.BrandID
	dbmodel.CategoryID = model.CategoryID
	dbmodel.ShortDescription = model.ShortDescription
	dbmodel.Description = model.Description
}

type ProductOutPutModel struct {
	ID           *uuid.UUID `gorm:"id"            json:"id"`
	Name         string     `gorm:"name"          json:"name"`
	Code         string     `gorm:"code"          json:"code"`
	BrandID      uuid.UUID  `gorm:"brand_id"      json:"brandId"`
	BrandName    string     `gorm:"brand_name"    json:"brandName"`
	CategoryID   uuid.UUID  `gorm:"category_id"   json:"categoryId"`
	CategoryName string     `gorm:"category_name" json:"categoryName"`
	Price        float32    `gorm:"price"         json:"price"`
	ItemID       uuid.UUID  `gorm:"item_id"       json:"itemId"`
}
