package dbseed

import (
	"time"

	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func seedColor(dbConn *gorm.DB) error {
	items := []models.Color{
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse("a15f0106-ee7d-424d-81d5-ceafcf86ef8d")
					return &id
				}(),
			},
			Name:     "آبی",
			Code:     "1",
			ColorHex: "#2762EB",
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse("adffee49-0e8d-4a3c-8318-53721931cda9")
					return &id
				}(),
			},
			Name:     "طوسی",
			Code:     "2",
			ColorHex: "#C7C7D1",
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse("db42920b-1a99-43f1-bb3e-b19ea294c6c7")
					return &id
				}(),
			},
			Name:     "سفید",
			Code:     "3",
			ColorHex: "#FFFFFF",
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse("a26cabb9-86ad-476b-8ffb-b3f83f013340")
					return &id
				}(),
			},
			Name:     "قرمز",
			Code:     "4",
			ColorHex: "#FF0000",
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse("48c9bd27-a51d-41b7-9011-a297e47ba5cb")
					return &id
				}(),
			},
			Name:     "سبز",
			Code:     "5",
			ColorHex: "#2ab57d",
		},
		{
			Model: models.Model{
				ID: func() *uuid.UUID {
					id := uuid.MustParse("a07a7c66-2b73-41c0-918c-16c73ce48e1e")
					return &id
				}(),
			},
			Name:     "نارنجی",
			Code:     "6",
			ColorHex: "#DD9654",
		},
	}

	for _, item := range items {

		var old models.Color
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
