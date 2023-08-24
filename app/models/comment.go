package models

import (
	"time"

	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/validations"
	dbmodels "github.com/esmailemami/eshop/models"
	datatypes "github.com/esmailemami/eshop/models/data_types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type CommentReqModel struct {
	Text           string                `json:"text"`
	Rate           int                   `json:"rate"`
	StrengthPoints datatypes.StringArray `json:"strengthPoints"`
	WeakPonits     datatypes.StringArray `json:"weakPonits"`
	ProductID      uuid.UUID             `json:"productId"`
}

func (model CommentReqModel) ValidateCreate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Text,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Rate,
			validation.Required.Error(consts.Required),
		),
		validation.Field(&model.ProductID,
			validation.Required.Error(consts.Required),
			validation.By(validations.ExistsInDB(&dbmodels.Product{}, "id", consts.ModelProductNotFound)),
		),
	)
}

func (model CommentReqModel) ValidateUpdate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Text,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Rate,
			validation.Required.Error(consts.Required),
			validation.By(validations.NumericValidator()),
		),
		validation.Field(&model.ProductID,
			validation.Required.Error(consts.Required),
			validation.By(validations.ExistsInDB(&dbmodels.Product{}, "id", consts.ModelProductNotFound)),
		),
	)
}

func (model CommentReqModel) ToDBModel() *dbmodels.Comment {
	return &dbmodels.Comment{
		Model: dbmodels.Model{
			ID: dbmodels.NewID(),
		},
		Text:           model.Text,
		Rate:           model.Rate,
		ProductID:      model.ProductID,
		StrengthPoints: model.StrengthPoints,
		WeakPonits:     model.WeakPonits,
		Status:         dbmodels.CommentStatusPending,
	}
}

func (model *CommentReqModel) MergeWithDBData(dbmodel *dbmodels.Comment) {
	dbmodel.Text = model.Text
	dbmodel.Rate = model.Rate
	dbmodel.StrengthPoints = model.StrengthPoints
	dbmodel.WeakPonits = model.WeakPonits
}

type ProductCommentOutPutModel struct {
	ID             *uuid.UUID            `gorm:"column:id"              json:"id"`
	CreatedAt      time.Time             `gorm:"column:created_at"      json:"createdAt"`
	UpdatedAt      time.Time             `gorm:"column:updated_at"      json:"updatedAt"`
	Text           string                `gorm:"column:text"            json:"text"`
	Rate           int                   `gorm:"column:rate"            json:"rate"`
	StrengthPoints datatypes.StringArray `gorm:"column:strength_points" json:"strengthPoints"`
	WeakPonits     datatypes.StringArray `gorm:"column:weak_ponits"     json:"weakPonits"`
	Username       string                `gorm:"column:username"        json:"username"`
}

type AdminCommentOutPutModel struct {
	ID             *uuid.UUID             `gorm:"column:id"              json:"id"`
	CreatedAt      time.Time              `gorm:"column:created_at"      json:"createdAt"`
	UpdatedAt      time.Time              `gorm:"column:updated_at"      json:"updatedAt"`
	Text           string                 `gorm:"column:text"            json:"text"`
	Rate           int                    `gorm:"column:rate"            json:"rate"`
	StrengthPoints datatypes.StringArray  `gorm:"column:strength_points" json:"strengthPoints"`
	WeakPonits     datatypes.StringArray  `gorm:"column:weak_ponits"     json:"weakPonits"`
	Username       string                 `gorm:"column:username"        json:"username"`
	ProductID      uuid.UUID              `gorm:"product_id"             json:"productId"`
	ProductName    string                 `gorm:"product_name"           json:"productName"`
	CommentStatus  dbmodels.CommentStatus `gorm:"comment_status"         json:"commentStatus"`
	AdminNote      *string                `gorm:"admin_note"                            json:"adminNote"`
}

type UserCommentOutPutModel struct {
	ID             *uuid.UUID             `gorm:"column:id"              json:"id"`
	CreatedAt      time.Time              `gorm:"column:created_at"      json:"createdAt"`
	UpdatedAt      time.Time              `gorm:"column:updated_at"      json:"updatedAt"`
	Text           string                 `gorm:"column:text"            json:"text"`
	Rate           int                    `gorm:"column:rate"            json:"rate"`
	StrengthPoints datatypes.StringArray  `gorm:"column:strength_points" json:"strengthPoints"`
	WeakPonits     datatypes.StringArray  `gorm:"column:weak_ponits"     json:"weakPonits"`
	ProductID      uuid.UUID              `gorm:"product_id"             json:"productId"`
	ProductName    string                 `gorm:"product_name"           json:"productName"`
	CommentStatus  dbmodels.CommentStatus `gorm:"comment_status"         json:"commentStatus"`
	AdminNote      *string                `gorm:"admin_note"                            json:"adminNote"`
}

type ChangeCommentStatus struct {
	Status dbmodels.CommentStatus `json:"status"`
	Note   *string                `json:"note,omitempty"`
}

func (model *ChangeCommentStatus) MergeWithDBData(dbmodel *dbmodels.Comment) {
	dbmodel.Status = model.Status
	dbmodel.AdminNote = model.Note
}
