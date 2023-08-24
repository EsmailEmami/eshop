package models

import (
	"time"

	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/validations"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type BrandReqModel struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (model BrandReqModel) ValidateCreate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Code,
			validation.Required.Error(consts.Required),
			validation.By(validations.Code()),
			validation.By(validations.NotExistsInDB(&dbmodels.Brand{}, "code", consts.ExistedCode)),
		),
	)
}

func (model BrandReqModel) ValidateUpdate(id uuid.UUID) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Code,
			validation.Required.Error(consts.Required),
			validation.By(validations.Code()),
			validation.By(validations.NotExistsInDBWithID(&dbmodels.Brand{}, "code", id, consts.ExistedCode)),
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
