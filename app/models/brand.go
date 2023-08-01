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

type BrandReqModel struct {
	Name   string    `json:"name"`
	Code   string    `json:"code"`
	FileID uuid.UUID `json:"fileId,omitempty"`
}

func (model BrandReqModel) ValidateCreate(db *gorm.DB) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Code,
			validation.Required.Error(consts.Required),
			validation.By(validations.Code()),
			validation.By(func(value interface{}) error {

				if dbpkg.Exists(db, &dbmodels.Brand{}, "code=?", value) {
					return errors.New(consts.ExistedCode)
				}

				return nil
			}),
		),
	)
}

func (model BrandReqModel) ValidateUpdate(db *gorm.DB) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Code,
			validation.Required.Error(consts.Required),
			validation.By(validations.Code()),
		),
		validation.Field(&model.FileID,
			validation.Required.Error(consts.Required),
			validation.By(func(value interface{}) error {
				if !dbpkg.Exists(db, &dbmodels.File{}, "id", value) {
					return errors.New("فایل یافت نشد")
				}

				return nil
			}),
		),
	)
}

func (model BrandReqModel) ToDBModel() *dbmodels.Brand {
	return &dbmodels.Brand{
		Model: dbmodels.Model{
			ID: dbmodels.NewID(),
		},
		Name:   model.Name,
		Code:   model.Code,
		FileID: uuid.MustParse(consts.FILE_DEFAULT_ID),
	}
}

func (model BrandReqModel) MergeWithDBData(dbmodel *dbmodels.Brand) {
	dbmodel.Name = model.Name
	dbmodel.Code = model.Code
}

type BrandOutPutModel struct {
	ID        *uuid.UUID        `gorm:"column:id"         json:"id"`
	CreatedAt time.Time         `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time         `gorm:"column:updated_at" json:"updatedAt"`
	Name      string            `gorm:"column:name"       json:"name"`
	Code      string            `gorm:"column:code"       json:"code"`
	FileID    uuid.UUID         `gorm:"column:file_id"    json:"fileId"`
	FileUrl   string            `gorm:"-"                 json:"fileUrl"`
	FileName  string            `gorm:"column:file_name"  json:"fileName"`
	FileType  dbmodels.FileType `gorm:"column:file_type"  json:"fileType"`
}
