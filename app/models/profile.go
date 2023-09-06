package models

import (
	"time"

	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/validations"
	dbmodels "github.com/esmailemami/eshop/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type UserDashboardInfoOutPutModel struct {
	ID        *uuid.UUID `gorm:"column:id"                              json:"id"`
	Username  string     `gorm:"column:username"                        json:"username"`
	FirstName *string    `gorm:"column:first_name"                      json:"firstName"`
	LastName  *string    `gorm:"column:last_name"                       json:"lastName"`
	Mobile    *string    `gorm:"column:mobile"                          json:"mobile"`
	RoleName  string     `gorm:"column:role_name"                       json:"roleName"`
	Email     *string    `gorm:"column:email"                           json:"email"`
}

type UserOrderOutPutModel struct {
	ID            *uuid.UUID             `gorm:"column:id"                         json:"id"`
	CreatedAt     time.Time              `gorm:"column:created_at"                 json:"createdAt"`
	Status        dbmodels.OrderStatus   `gorm:"column:status"                     json:"status"`
	Price         float64                `gorm:"column:price"                      json:"price"`
	Code          string                 `gorm:"column:code"                       json:"code"`
	PaidAt        *time.Time             `gorm:"column:paid_at"                    json:"paidAt"`
	TotalPrice    float64                `gorm:"column:total_price"                json:"totalPrice"`
	DiscountPrice *float64               `gorm:"column:discount_price"             json:"discountPrice,omitempty"`
	DiscountValue *float64               `gorm:"column:discount_value"             json:"discountValue,omitempty"`
	DiscountType  *dbmodels.DiscountType `gorm:"column:discount_type"              json:"discountType,omitempty"`
	FileUrls      []string               `gorm:"-"                                 json:"fileUrls"`
}

type UserProfileUpdateModel struct {
	Username  string  `json:"username"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Mobile    *string `json:"mobile"`
	Email     *string `json:"email"`
}

func (model UserProfileUpdateModel) ValidateUpdate(id uuid.UUID) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(&model.Username,
			validation.Required.Error(consts.Required),
			validation.By(validations.UserName()),
			validation.By(validations.NotExistsInDBWithID(&dbmodels.User{}, "username", id, consts.UsernameAlreadyExists)),
		),
		validation.Field(&model.FirstName,
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.LastName,
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Mobile,
			validation.By(validations.IsValidMobileNumber()),
		),
	)
}

func (model *UserProfileUpdateModel) MergeWithDBData(dbmodel *dbmodels.User) {
	dbmodel.Username = model.Username
	dbmodel.FirstName = model.FirstName
	dbmodel.LastName = model.LastName
	dbmodel.Mobile = model.Mobile
	dbmodel.Email = model.Email
}
