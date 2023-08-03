package controllers

import (
	"net/http"

	"github.com/esmailemami/eshop/app"
	appmodels "github.com/esmailemami/eshop/app/models"
	"github.com/esmailemami/eshop/consts"
	"github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/errors"
	"github.com/esmailemami/eshop/models"
)

// GetUser godoc
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} appmodels.UserDashboardInfoOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user [get]
func GetUser(ctx *app.HttpContext) error {
	user, err := ctx.GetUser()
	if err != nil {
		return errors.NewUnauthorizedError(consts.UnauthorizedError, err)
	}

	data := appmodels.UserDashboardInfoOutPutModel{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Mobile:    user.Mobile,
	}

	if user.Role != nil {
		data.RoleName = user.Role.Name
	}

	return ctx.JSON(data, http.StatusOK)
}

// GetUser godoc
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param status  query  models.OrderStatus  false  "Order Status"
// @Success 200 {object} []appmodels.UserOrderOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/orders [get]
func GetUserOrders(ctx *app.HttpContext) error {
	user, err := ctx.GetUser()
	if err != nil {
		return errors.NewUnauthorizedError(consts.UnauthorizedError, err)
	}

	baseDB := db.MustGormDBConn(ctx)
	orderDB := baseDB.Model(&models.Order{}).Where("created_by_id=?", *user.ID)

	orderStatus, ok := ctx.GetParam("status")
	if ok {
		orderDB = orderDB.Where("status=?", orderStatus)
	}

	var orders []appmodels.UserOrderOutPutModel

	if err := orderDB.Find(&orders).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	// files

	for i, order := range orders {
		files := []struct {
			FileType       models.FileType `gorm:"column:file_type"`
			UniqueFileName string          `gorm:"column:unique_file_name"`
		}{}

		if err := baseDB.Table("order_item oi").
			Joins("INNER JOIN product_item pi ON pi.id = oi.product_item_id").
			Joins("INNER JOIN product p ON p.id = pi.product_id").
			Joins("CROSS JOIN LATERAL (?) as pf", baseDB.Table("product_file_map pf").
				Select("file_id").
				Where("pf.product_id = p.id").Limit(1),
			).
			Joins("INNER JOIN file f ON f.id = pf.file_id").
			Where("p.deleted_at IS NULL AND oi.order_id = ?", *order.ID).
			Select("f.file_type, unique_file_name").
			Find(&files).Error; err != nil {
			return errors.NewInternalServerError(consts.InternalServerError, err)
		}

		for _, file := range files {
			order.FileUrls = append(order.FileUrls, models.GetFileUrl(file.FileType, file.UniqueFileName))
		}

		orders[i] = order
	}

	return ctx.JSON(orders, http.StatusOK)
}
