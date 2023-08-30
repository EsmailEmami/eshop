package models

import (
	"time"

	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/validations"
	"github.com/esmailemami/eshop/models"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type AppPicReqModel struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Priority    int                 `json:"priority"`
	AppPicType  dbmodels.AppPicType `json:"appPicType"`
	Url         string              `json:"url"`
}

func (model AppPicReqModel) ValidateCreate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Title,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Description,
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Url,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Priority,
			validation.Required.Error(consts.Required),
			validation.By(validations.NotExistsInDBWithCond(&models.AppPic{}, "priority", consts.OrderAlreadyRegistered,
				"app_pic_type=?", model.AppPicType)),
		),
	)
}

func (model AppPicReqModel) ValidateUpdate(id uuid.UUID) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Title,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Description,
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Url,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Priority,
			validation.Required.Error(consts.Required),
			validation.By(validations.NotExistsInDBWithCond(&models.AppPic{}, "priority", consts.OrderAlreadyRegistered,
				"id != ? AND app_pic_type=?", id, model.AppPicType)),
		),
	)
}

func (model *AppPicReqModel) ToDBModel() *dbmodels.AppPic {
	return &dbmodels.AppPic{
		Model: dbmodels.Model{
			ID: dbmodels.NewID(),
		},
		Title:       model.Title,
		Description: model.Description,
		FileID:      uuid.MustParse(consts.FILE_DEFAULT_ID),
		Priority:    model.Priority,
		AppPicType:  model.AppPicType,
		Url:         model.Url,
	}
}

func (model *AppPicReqModel) MergeWithDBData(dbmodel *dbmodels.AppPic) {
	dbmodel.Title = model.Title
	dbmodel.Description = model.Description
	dbmodel.Priority = model.Priority
	dbmodel.AppPicType = model.AppPicType
	dbmodel.Url = model.Url
}

type AppPicOutPutModel struct {
	ID          *uuid.UUID          `gorm:"column:id"                json:"id"`
	CreatedAt   time.Time           `gorm:"column:created_at"        json:"createdAt"`
	UpdatedAt   time.Time           `gorm:"column:updated_at"        json:"updatedAt"`
	Priority    int                 `gorm:"column:priority"          json:"priority"`
	AppPicType  dbmodels.AppPicType `gorm:"column:app_pic_type"      json:"appPicType"`
	FileID      uuid.UUID           `gorm:"column:file_id"           json:"fileId"`
	Title       string              `gorm:"column:title"             json:"title"`
	Description string              `gorm:"column:description"       json:"description"`
	Url         string              `gorm:"column:url"               json:"url"`
	FileType    dbmodels.FileType   `gorm:"column:file_type"         json:"fileType"`
	FileName    string              `gorm:"column:file_name"         json:"fileName"`
	FileUrl     string              `gorm:"-"                        json:"fileUrl"`
}
