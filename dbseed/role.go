package dbseed

import (
	"time"

	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func seedRole(dbConn *gorm.DB) error {
	items := []models.Role{
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse(consts.ROLE_ROOT_ID)
					return &id
				}(),
			},
			Name:     "root",
			Code:     "1",
			IsSystem: true,
			Permissions: models.RolePermissions{
				models.ACTION_CAN_LOGIN_ADMIN,
				models.ACTION_CAN_LOGIN_USER,
				models.ACTION_ROLE_ADMIN_LIST,
				models.ACTION_ROLE_ADMIN_INFO,
				models.ACTION_ROLE_ADMIN_CREATE,
				models.ACTION_ROLE_ADMIN_UPDATE,
				models.ACTION_ROLE_ADMIN_DELETE,
				models.ACTION_ROLE_ADMIN_PERMISSIONS,
			},
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse(consts.ROLE_ADMIN_ID)
					return &id
				}(),
			},
			Name:     "admin",
			Code:     "2",
			IsSystem: true,
			Permissions: models.RolePermissions{
				models.ACTION_CAN_LOGIN_ADMIN,
				models.ACTION_CAN_LOGIN_USER,
				models.ACTION_ROLE_ADMIN_LIST,
				models.ACTION_ROLE_ADMIN_INFO,
				models.ACTION_ROLE_ADMIN_PERMISSIONS,
			},
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse(consts.ROLE_NO_ACCESS_ID)
					return &id
				}(),
			},
			Name:     "no access",
			Code:     "3",
			IsSystem: true,
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse(consts.ROLE_USER)
					return &id
				}(),
			},
			Name:     "user",
			Code:     "4",
			IsSystem: true,
			Permissions: models.RolePermissions{
				models.ACTION_CAN_LOGIN_USER,
			},
		},
	}

	for _, item := range items {

		var old models.Role
		err := dbConn.Where("id", item.ID).First(&old).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}

		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		if item.CreatedAt.IsZero() {
			item.CreatedAt = time.Now()
		}
		err = dbConn.Save(&item).Error
		if err != nil {
			return err
		}
	}

	return nil
}
