package dbseed

import (
	"time"

	"github.com/esmailemami/eshop/consts"
	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func seedBrand(dbConn *gorm.DB) error {
	items := []models.Brand{
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse("95231911-ba44-47b1-aac4-47a544546ebd")
					return &id
				}(),
			},
			Name:   "اپل",
			Code:   "1",
			FileID: uuid.MustParse(consts.FILE_BRAND_APPLE_ID),
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse("5ea2756a-4647-4002-95a2-f4304b2f884d")
					return &id
				}(),
			},
			Name:   "هواوی",
			Code:   "2",
			FileID: uuid.MustParse(consts.FILE_BRAND_HUAWEI_ID),
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse("01b556ea-72a6-4f9b-adf7-889a64073748")
					return &id
				}(),
			},
			Name:   "شیائومی",
			Code:   "3",
			FileID: uuid.MustParse(consts.FILE_BRAND_XIAOMI_ID),
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse("441af1b3-8fea-4209-8ae4-4f4243d5dcca")
					return &id
				}(),
			},
			Name:   "سامسونگ",
			Code:   "4",
			FileID: uuid.MustParse(consts.FILE_BRAND_SAMSUNG_ID),
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse("cc96d01e-72c5-4114-a47b-5879d931295a")
					return &id
				}(),
			},
			Name:   "کوکا کولا",
			Code:   "5",
			FileID: uuid.MustParse(consts.FILE_BRAND_COCACOLA_ID),
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse("1fa66789-be01-4192-b5a8-5535867944f3")
					return &id
				}(),
			},
			Name:   "نستله",
			Code:   "6",
			FileID: uuid.MustParse(consts.FILE_BRAND_NESTELE_ID),
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse("3c84acd8-3c0b-4aa2-b0da-c3630d7b2370")
					return &id
				}(),
			},
			Name:   "پپسی",
			Code:   "7",
			FileID: uuid.MustParse(consts.FILE_BRAND_PEPSI_ID),
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse("d4fb47f3-50f5-43f0-95f9-ccdbc84f2a32")
					return &id
				}(),
			},
			Name:   "اوپو",
			Code:   "8",
			FileID: uuid.MustParse(consts.FILE_BRAND_OPPO_ID),
		},
	}

	for _, item := range items {

		var old models.Brand
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
