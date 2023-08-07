package models

import (
	"errors"

	"github.com/esmailemami/eshop/app/validations"
	"github.com/esmailemami/eshop/consts"
	dbpkg "github.com/esmailemami/eshop/db"
	dbmodels "github.com/esmailemami/eshop/models"
	datatypes "github.com/esmailemami/eshop/models/data_types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductReqModel struct {
	Name             string                `json:"name"`
	Code             string                `json:"code"`
	BrandID          uuid.UUID             `json:"brandId"`
	CategoryID       uuid.UUID             `json:"categoryId"`
	Description      string                `json:"description"`
	ShortDescription string                `json:"shortDescription"`
	TopFeatures      datatypes.StringArray `json:"topFeatures"`
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
					return errors.New(consts.ModelBrandNotFound)
				}

				return nil
			}),
		),
		validation.Field(&model.CategoryID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {

				if !dbpkg.Exists(db, &dbmodels.Category{}, "id=?", value) {
					return errors.New(consts.ModelCategoryNotFound)
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
					return errors.New(consts.ModelBrandNotFound)
				}

				return nil
			}),
		),
		validation.Field(&model.CategoryID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {

				if !dbpkg.Exists(db, &dbmodels.Category{}, "id=?", value) {
					return errors.New(consts.ModelCategoryNotFound)
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
		TopFeatures:      model.TopFeatures,
	}
}

func (model *ProductReqModel) MergeWithDBData(dbmodel *dbmodels.Product) {
	dbmodel.Name = model.Name
	dbmodel.Code = model.Code
	dbmodel.BrandID = model.BrandID
	dbmodel.CategoryID = model.CategoryID
	dbmodel.ShortDescription = model.ShortDescription
	dbmodel.Description = model.Description
	dbmodel.TopFeatures = model.TopFeatures
}

type ProductWithItemOutPutModel struct {
	ID           *uuid.UUID        `gorm:"column:id"               json:"id"`
	Name         string            `gorm:"column:name"             json:"name"`
	Code         string            `gorm:"column:code"             json:"code"`
	BrandID      uuid.UUID         `gorm:"column:brand_id"         json:"brandId"`
	BrandName    string            `gorm:"column:brand_name"       json:"brandName"`
	CategoryID   uuid.UUID         `gorm:"column:category_id"      json:"categoryId"`
	CategoryName string            `gorm:"column:category_name"    json:"categoryName"`
	Price        float32           `gorm:"column:price"            json:"price"`
	ItemID       uuid.UUID         `gorm:"column:item_id"          json:"itemId"`
	FileType     dbmodels.FileType `gorm:"column:file_type" json:"-"`
	FileName     string            `gorm:"column:file_name"        json:"-"`
	FileUrl      string            `gorm:"column:file_url"         json:"fileUrl"`
}

type ProductOutPutModel struct {
	ID            *uuid.UUID        `gorm:"column:id"               json:"id"`
	Name          string            `gorm:"column:name"             json:"name"`
	Code          string            `gorm:"column:code"             json:"code"`
	BrandID       uuid.UUID         `gorm:"column:brand_id"         json:"brandId"`
	BrandName     string            `gorm:"column:brand_name"       json:"brandName"`
	CategoryID    uuid.UUID         `gorm:"column:category_id"      json:"categoryId"`
	CategoryName  string            `gorm:"column:category_name"    json:"categoryName"`
	BrandFileType dbmodels.FileType `gorm:"column:brand_file_type"  json:"-"`
	BrandFileName string            `gorm:"column:brand_file_name"  json:"-"`
	BrandFileUrl  string            `gorm:"column:brand_file_url"   json:"brandFileUrl"`
}

// ------------------- suggestion products -----------------------

type SuggestionProductOutPutModel struct {
	ProductID     *uuid.UUID                        `gorm:"column:product_id"                json:"productId"`
	Name          string                            `gorm:"column:name"                      json:"name"`
	Colors        []ProductItemInfoColorOutPutModel `gorm:"-"                                json:"colors"`
	Files         []ProductItemFileOutPutModel      `gorm:"-"                                json:"files"`
	ProductItemID uuid.UUID                         `gorm:"product_item_id"                  json:"productItemId"`
	ColorID       uuid.UUID                         `gorm:"color_id"                         json:"colorId"`
	TopFeatures   datatypes.StringArray             `gorm:"column:top_features"              json:"topFeatures"`
}

type ProductAdminOutPutModel struct {
	ID               *uuid.UUID            `gorm:"column:id"                              json:"id"`
	Name             string                `gorm:"column:name"                            json:"name"`
	Code             string                `gorm:"column:code"                            json:"code"`
	BrandID          uuid.UUID             `gorm:"column:brand_id"         json:"brandId"`
	BrandName        string                `gorm:"column:brand_name"       json:"brandName"`
	CategoryID       uuid.UUID             `gorm:"column:category_id"      json:"categoryId"`
	CategoryName     string                `gorm:"column:category_name"    json:"categoryName"`
	Description      string                `gorm:"column:description"                     json:"description"`
	ShortDescription string                `gorm:"column:short_description"               json:"shortDescription"`
	TopFeatures      datatypes.StringArray `gorm:"column:top_features"                    json:"topFeatures"`
}
