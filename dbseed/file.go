package dbseed

import (
	"time"

	"github.com/esmailemami/eshop/consts"
	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func seedFile(dbConn *gorm.DB) error {
	items := []models.File{
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse(consts.DEFAULT_FILE_ID)
					return &id
				}(),
			},
			MimeType:       "png",
			OriginalName:   "default picture",
			UniqueFileName: consts.DEFAULT_FILE_ID + ".png",
			FileType:       models.FileTypeSystematic,
		},
	}

	for _, item := range items {

		var old models.File
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
