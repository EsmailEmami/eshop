package user

import (
	"context"

	dbpkg "github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/models"
)

func UserExistsWithUsername(value string) (bool, error) {
	db := dbpkg.MustGormDBConn(context.Background())

	var count int64
	err := db.Model(&models.User{}).Where("username = ?", value, value).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}

	return true, nil
}

func GetUserByUsername(value string) (*models.User, error) {
	db := dbpkg.MustGormDBConn(context.Background())

	var user models.User
	err := db.Model(&models.User{}).Preload("Role").Where("username = ?", value).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
