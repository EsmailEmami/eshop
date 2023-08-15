package models

import (
	"context"
	"errors"
	"time"

	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/validations"
	"github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginInputModel struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	IP        string `json:"-"`
	UserAgent string `json:"-"`
}

func (model LoginInputModel) Validate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(
			&model.Username,
			validation.Required.Error(consts.Required),
			validation.By(validations.UserName()),
		),
		validation.Field(
			&model.Password,
			validation.Required.Error(consts.Required),
			validation.By(validations.StrongPassword()),
		),
	)
}

type LoginOutputModel struct {
	TokenID   uuid.UUID            `json:"-"`
	Token     string               `json:"token"`
	ExpiresAt time.Time            `json:"expiresAt"`
	ExpiresIn int64                `json:"expiresIn"`
	User      LoginOutputUserModel `json:"user"`
}

type LoginOutputUserModel struct {
	Username  string  `json:"username"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}

type RegisterInputModel struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

func (model RegisterInputModel) Validate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(
			&model.Username,
			validation.Required.Error(consts.Required),
			validation.By(validations.UserName()),
			validation.By(func(value interface{}) error {
				if db.Exists(
					db.MustGormDBConn(context.Background()),
					&models.User{},
					"username = ?",
					value,
				) {
					return errors.New(consts.UsernameAlreadyExists)
				}

				return nil
			}),
		),
		validation.Field(
			&model.Password,
			validation.Required.Error(consts.Required),
			validation.By(validations.StrongPassword()),
			validation.By(func(value interface{}) error {
				if value.(string) != model.PasswordConfirmation {
					return errors.New(consts.PasswordMismatch)
				}
				return nil
			}),
		),
	)
}

func (model RegisterInputModel) ToDBModel() *models.User {
	pass, _ := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)

	return &models.User{
		Model: models.Model{
			ID: func() *uuid.UUID {
				id := uuid.New()
				return &id
			}(),
		},
		Username: model.Username,
		Password: string(pass),
	}
}

type RecoveryPasswordModel struct {
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

func (model RecoveryPasswordModel) Validate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(
			&model.Password,
			validation.Required.Error(consts.Required),
			validation.By(validations.StrongPassword()),
			validation.By(func(value interface{}) error {
				if value.(string) != model.PasswordConfirmation {
					return errors.New(consts.PasswordMismatch)
				}
				return nil
			}),
		),
	)
}

type RecoveryPasswordReqModel struct {
	PhoneNumberOrEmailAddress string `json:"phoneNumberOrEmailAddress"`
}

func (model RecoveryPasswordReqModel) Validate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(
			&model.PhoneNumberOrEmailAddress,
			validation.Required.Error(consts.Required),
		),
	)
}
