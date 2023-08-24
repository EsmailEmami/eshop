package models

import (
	"time"

	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/validations"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type RoleReqModel struct {
	Name        string                   `json:"name"`
	Code        string                   `json:"code"`
	IsSystem    bool                     `json:"isSystem"`
	Permissions dbmodels.RolePermissions `json:"permissions"`
}

func (model RoleReqModel) ValidateCreate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.IsSystem,
			validation.Required.Error(consts.Required),
		),
		validation.Field(&model.Code,
			validation.Required.Error(consts.Required),
			validation.By(validations.Code()),
			validation.By(validations.NotExistsInDB(&dbmodels.Role{}, "code", consts.ExistedCode)),
		),
	)
}

func (model RoleReqModel) ValidateUpdate(id uuid.UUID) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.IsSystem,
			validation.Required.Error(consts.Required),
		),
		validation.Field(&model.Code,
			validation.Required.Error(consts.Required),
			validation.By(validations.Code()),
			validation.By(validations.NotExistsInDBWithID(&dbmodels.Role{}, "code", id, consts.ExistedCode)),
		),
	)
}

func (model *RoleReqModel) ToDBModel() *dbmodels.Role {
	return &dbmodels.Role{
		Model: dbmodels.Model{
			ID: dbmodels.NewID(),
		},
		Name:        model.Name,
		Code:        model.Code,
		IsSystem:    model.IsSystem,
		Permissions: model.Permissions,
	}
}

func (model *RoleReqModel) MergeWithDBData(dbmodel *dbmodels.Role) {
	dbmodel.Name = model.Name
	dbmodel.Code = model.Code
	dbmodel.IsSystem = model.IsSystem
	dbmodel.Permissions = model.Permissions
}

type RoleOutPutModel struct {
	ID          *uuid.UUID               `gorm:"column:id"          json:"id"`
	CreatedAt   time.Time                `gorm:"column:created_at"  json:"createdAt"`
	UpdatedAt   time.Time                `gorm:"column:updated_at"  json:"updatedAt"`
	Name        string                   `gorm:"column:name"        json:"name"`
	Code        string                   `gorm:"column:code"        json:"code"`
	IsSystem    bool                     `gorm:"column:is_system"   json:"isSystem"`
	Permissions dbmodels.RolePermissions `gorm:"column:permissions" json:"permissions"`
}
