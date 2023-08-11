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

type RoleReqModel struct {
	Name        string                   `json:"name"`
	Code        string                   `json:"code"`
	IsSystem    bool                     `json:"isSystem"`
	Permissions dbmodels.RolePermissions `json:"permissions"`
}

func (model RoleReqModel) ValidateCreate(db *gorm.DB) error {
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
			validation.By(func(value interface{}) error {

				if dbpkg.Exists(db, &dbmodels.Role{}, "code=?", value) {
					return errors.New(consts.ExistedCode)
				}

				return nil
			}),
		),
	)
}

func (model RoleReqModel) ValidateUpdate() error {
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
	ID          *uuid.UUID               `gorm:"column:id"         json:"id"`
	CreatedAt   time.Time                `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time                `gorm:"column:updated_at" json:"updatedAt"`
	Name        string                   `gorm:"column:name"       json:"name"`
	Code        string                   `gorm:"column:code"       json:"code"`
	Permissions dbmodels.RolePermissions `gorm:"column:permissions" json:"permissions"`
}
