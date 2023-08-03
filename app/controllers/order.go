package controllers

import (
	"net/http"

	"github.com/esmailemami/eshop/app"
	appmodels "github.com/esmailemami/eshop/app/models"
	"github.com/esmailemami/eshop/consts"
	"github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/errors"
)

// GetOrders godoc
// @Tags Orders
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} appmodels.OrderOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /order [get]
func GetOrder(ctx *app.HttpContext) error {
	user, err := ctx.GetUser()
	if err != nil {
		return errors.NewUnauthorizedError(err.Error(), err)
	}

	baseDB := db.MustGormDBConn(ctx)

	data := appmodels.OrderOutPutModel{
		Items: []appmodels.OrderItemOutPutModel{},
	}
	var orderItems []appmodels.OrderItemOutPutModel

	if err := baseDB.Table("order_item oi").
		Joins(`INNER JOIN "order" o ON o.id = oi.order_id`).
		Joins("INNER JOIN product_item pi2 ON pi2.id = oi.product_item_id").
		Joins("INNER JOIN product p ON	p.id = pi2.product_id").
		Joins("CROSS JOIN LATERAL (?) as pf", baseDB.Table("product_file_map pfm").
			Select("file_id").
			Where("pfm.product_id = p.id").
			Limit(1),
		).
		Joins("INNER JOIN file f ON f.id = pf.file_id").
		Where("o.status = 0 AND o.created_by_id = ?", *user.ID).
		Select(`oi.id, oi.product_item_id,p."name" AS product_name,pi2.price,oi.quantity,f.file_type, f.unique_file_name AS file_name`).
		Find(&orderItems).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	for _, orderItem := range orderItems {
		orderItem.FileUrl = orderItem.FileType.GetDirectory() + "/" + orderItem.FileName
		orderItem.TotalPrice = orderItem.Price * float64(orderItem.Quantity)
		data.Items = append(data.Items, orderItem)
		data.Price += orderItem.Price
	}

	return ctx.JSON(data, http.StatusOK)
}

// GetOrders godoc
// @Tags Orders
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} appmodels.OrderOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /order [get]
