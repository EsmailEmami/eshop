package dbseed

import (
	"time"

	"github.com/esmailemami/eshop/consts"
	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func seedUsers(dbConn *gorm.DB) error {
	items := []models.User{
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse("68633fb0-b9a4-4d2f-a441-73faf5e3fa15")
					return &id
				}(),
			},
			FirstName: "محمد",
			LastName:  "غلامی",
			Username:  "mgh",
			Password:  "$2a$10$LxCfacwwZ2TYmJJRBU0RGu/nY.15kiqPaQq8IAWBhjxSXbCTAep7u",
			IsSystem:  true,
			Mobile:    "09171209668",
			Enabled:   true,
			NoeSematID: func() *uuid.UUID {
				id := uuid.MustParse(consts.NOE_SEMAT_ROOT_ID)
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
			FirstName: "عبدالله",
			LastName:  "محمدی",
			Username:  "abdollahm66",
			Password:  "$2a$12$2r6Z7XiJgUTQpLcxcGd15uUezXug4fmDYz4VV2jfHV7kUwZHnqGLS",
			IsSystem:  true,
			Mobile:    "09179818310",
			Enabled:   true,
			NoeSematID: func() *uuid.UUID {
				id := uuid.MustParse(consts.NOE_SEMAT_ROOT_ID)
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
				"first_name":   item.FirstName,
				"last_name":    item.LastName,
				"username":     item.Username,
				"is_system":    item.IsSystem,
				"mobile":       item.Mobile,
				"enabled":      item.Enabled,
				"noe_semat_id": item.NoeSematID.String(),
				"updated_at":   time.Now(),
			}).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}
