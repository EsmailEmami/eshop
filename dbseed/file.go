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
					id := uuid.MustParse(consts.FILE_DEFAULT_ID)
					return &id
				}(),
			},
			MimeType:       "image/jpeg",
			OriginalName:   "default picture",
			UniqueFileName: consts.FILE_DEFAULT_ID + ".jpg",
			FileType:       models.FileTypeSystematic,
			Extension:      "jpg",
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse(consts.FILE_BRAND_APPLE_ID)
					return &id
				}(),
			},
			MimeType:       "image/png",
			OriginalName:   "APPLE LOGO",
			UniqueFileName: consts.FILE_BRAND_APPLE_ID + ".png",
			FileType:       models.FileTypeBrand,
			Extension:      "png",
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse(consts.FILE_BRAND_HUAWEI_ID)
					return &id
				}(),
			},
			MimeType:       "image/png",
			OriginalName:   "HUAWEI LOGO",
			UniqueFileName: consts.FILE_BRAND_HUAWEI_ID + ".png",
			FileType:       models.FileTypeBrand,
			Extension:      "png",
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse(consts.FILE_BRAND_XIAOMI_ID)
					return &id
				}(),
			},
			MimeType:       "image/png",
			OriginalName:   "XIAOMI LOGO",
			UniqueFileName: consts.FILE_BRAND_XIAOMI_ID + ".png",
			FileType:       models.FileTypeBrand,
			Extension:      "png",
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse(consts.FILE_BRAND_SAMSUNG_ID)
					return &id
				}(),
			},
			MimeType:       "image/png",
			OriginalName:   "SAMSUNG LOGO",
			UniqueFileName: consts.FILE_BRAND_SAMSUNG_ID + ".png",
			FileType:       models.FileTypeBrand,
			Extension:      "png",
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse(consts.FILE_BRAND_COCACOLA_ID)
					return &id
				}(),
			},
			MimeType:       "image/png",
			OriginalName:   "COCACOLA LOGO",
			UniqueFileName: consts.FILE_BRAND_COCACOLA_ID + ".png",
			FileType:       models.FileTypeBrand,
			Extension:      "png",
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse(consts.FILE_BRAND_NESTELE_ID)
					return &id
				}(),
			},
			MimeType:       "image/png",
			OriginalName:   "NESTELE LOGO",
			UniqueFileName: consts.FILE_BRAND_NESTELE_ID + ".png",
			FileType:       models.FileTypeBrand,
			Extension:      "png",
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse(consts.FILE_BRAND_PEPSI_ID)
					return &id
				}(),
			},
			MimeType:       "image/png",
			OriginalName:   "PEPSI LOGO",
			UniqueFileName: consts.FILE_BRAND_PEPSI_ID + ".png",
			FileType:       models.FileTypeBrand,
			Extension:      "png",
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse(consts.FILE_BRAND_OPPO_ID)
					return &id
				}(),
			},
			MimeType:       "image/png",
			OriginalName:   "OPPO LOGO",
			UniqueFileName: consts.FILE_BRAND_OPPO_ID + ".png",
			FileType:       models.FileTypeBrand,
			Extension:      "png",
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
