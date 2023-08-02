package order

import (
	"math/rand"

	dbpkg "github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetOpenOrder(db *gorm.DB, userID uuid.UUID) (*models.Order, error) {
	var order models.Order

	if err := db.Model(&models.Order{}).Find(&order, "created_by_id=? AND status = 0", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			order = models.Order{
				Status: models.OrderStatusOpen,
			}

			// check the code is unique
			for {
				order.Code = generateOrderCode(9)

				if !dbpkg.Exists(db, &models.Order{}, "code=?", order.Code) {
					break
				}
			}

			if err := db.Create(&order).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &order, nil
}

func generateOrderCode(length int) string {
	charSet := "0123456789"

	code := make([]byte, length)
	for i := 0; i < length; i++ {
		code[i] = charSet[rand.Intn(len(charSet))]
	}

	return string(code)
}
