package models

import (
	"time"

	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/validations"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type UserReqModel struct {
	Username  string     `json:"username"`
	FirstName *string    `json:"firstName,omitempty"`
	LastName  *string    `json:"lastName,omitempty"`
	Mobile    *string    `json:"mobile,omitempty"`
	RoleID    *uuid.UUID `json:"roleId,omitempty"`
	Email     *string    `json:"email,omitempty"`
	IsSystem  bool       `json:"isSystem"`
	Enabled   bool       `json:"enabled"`
}

func (model UserReqModel) ValidateCreate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Username,
			validation.Required.Error(consts.Required),
			validation.By(validations.UserName()),
			validation.By(validations.NotExistsInDB(&dbmodels.User{}, "username", consts.UsernameAlreadyExists)),
		),
		validation.Field(&model.RoleID,
			validation.By(validations.ExistsInDB(&dbmodels.Role{}, "id", consts.ModelRoleNotFound)),
		),
	)
}

func (model UserReqModel) ValidateUpdate(id uuid.UUID) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Username,
			validation.Required.Error(consts.Required),
			validation.By(validations.UserName()),
			validation.By(validations.NotExistsInDBWithID(&dbmodels.User{}, "username", id, consts.UsernameAlreadyExists)),
		),
		validation.Field(&model.RoleID,
			validation.By(validations.ExistsInDB(&dbmodels.Role{}, "id", consts.ModelRoleNotFound)),
		),
	)
}

func (model *UserReqModel) ToDBModel() *dbmodels.User {
	return &dbmodels.User{
		Model: dbmodels.Model{
			ID: dbmodels.NewID(),
		},
		Username:  model.Username,
		FirstName: model.FirstName,
		LastName:  model.LastName,
		Mobile:    model.Mobile,
		RoleID:    model.RoleID,
		Email:     model.Email,
		IsSystem:  model.IsSystem,
		Enabled:   model.Enabled,
	}
}

func (model *UserReqModel) MergeWithDBData(dbmodel *dbmodels.User) {
	dbmodel.Username = model.Username
	dbmodel.FirstName = model.FirstName
	dbmodel.LastName = model.LastName
	dbmodel.Mobile = model.Mobile
	dbmodel.RoleID = model.RoleID
	dbmodel.Email = model.Email
	dbmodel.IsSystem = model.IsSystem
	dbmodel.Enabled = model.Enabled
}

type UserOutPutModel struct {
	ID        *uuid.UUID `gorm:"column:id"                              json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at"                      json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at"                      json:"updatedAt"`
	Username  string     `gorm:"column:username"                        json:"username"`
	Password  string     `gorm:"column:password"                        json:"-"`
	FirstName *string    `gorm:"column:first_name"                      json:"firstName"`
	LastName  *string    `gorm:"column:last_name"                       json:"lastName"`
	Mobile    *string    `gorm:"column:mobile"                          json:"mobile"`
	RoleID    *uuid.UUID `gorm:"column:role_id"                         json:"roleId"`
	Email     *string    `gorm:"email"                                  json:"email"`
	IsSystem  bool       `gorm:"column:is_system"                       json:"isSystem"`
	Enabled   bool       `gorm:"column:enabled"                         json:"enabled"`
}
