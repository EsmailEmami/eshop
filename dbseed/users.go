package dbseed

import (
	"time"

	"github.com/esmailemami/eshop/consts"
	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func seedUser(dbConn *gorm.DB) error {
	items := []models.User{
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse("68633fb0-b9a4-4d2f-a441-73faf5e3fa15")
					return &id
				}(),
			},
			FirstName: "اسماعیل",
			LastName:  "امامی",
			Username:  "esmailemami",
			Password:  "$2a$10$2oV2MylgwZftP47vL/ndteC6tzmcY85qRNo/5FTCeS403eL8zo9Yq",
			IsSystem:  true,
			Mobile:    "09903669556",
			Enabled:   true,
			Email:     "esmailemami84@gmail.com",
			RoleID: func() *uuid.UUID {
				id := uuid.MustParse(consts.ROLE_ROOT_ID)
				return &id
			}(),
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse("10dcab09-a39a-4238-9c27-314267be43a1")
					return &id
				}(),
			},
			FirstName: "علیرضا",
			LastName:  "صفری",
			Username:  "alireza83safari",
			Password:  "$2a$10$2oV2MylgwZftP47vL/ndteC6tzmcY85qRNo/5FTCeS403eL8zo9Yq",
			IsSystem:  true,
			Mobile:    "09903669556",
			Email:     "alireza83safarii@gmail.com",
			Enabled:   true,
			RoleID: func() *uuid.UUID {
				id := uuid.MustParse(consts.ROLE_ROOT_ID)
				return &id
			}(),
		},
	}

	for _, item := range items {

		var old models.User
		err := dbConn.Where("id", item.ID).First(&old).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}

		if err != nil && err == gorm.ErrRecordNotFound {
			err = dbConn.Save(&item).Error
			if err != nil {
				return err
			}
		} else {
			err = dbConn.Model(&models.User{}).Where("id", item.ID).UpdateColumns(map[string]any{
				"first_name": item.FirstName,
				"last_name":  item.LastName,
				"username":   item.Username,
				"is_system":  item.IsSystem,
				"mobile":     item.Mobile,
				"enabled":    item.Enabled,
				"role_id":    item.RoleID.String(),
				"updated_at": time.Now(),
				"email":      item.Email,
			}).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}
