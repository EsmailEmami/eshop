package controllers

import (
	"net/http"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/errors"
	appmodels "github.com/esmailemami/eshop/app/models"
	dbpkg "github.com/esmailemami/eshop/db"
)

// GetRevenueByCategory godoc
// @Tags Reports
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} []appmodels.RevenueByCategory
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/report/revenueByCategory [get]
func ReportRevenueByCategory(ctx *app.HttpContext) error {
	var (
		db                 = dbpkg.MustGormDBConn(ctx)
		data               = []appmodels.RevenueByCategory{}
		totalPrice float64 = 0
	)

	subquery := db.Table("order_item as oi").
		Select("product.category_id as category_id, sum(oi.price) as price").
		Joins("INNER JOIN \"order\" o ON o.id = oi.order_id").
		Joins("INNER JOIN product_item pi2 ON pi2.id = oi.product_item_id").
		Joins("INNER JOIN product ON product.id = pi2.product_id").
		Where("o.deleted_at IS NULL AND oi.deleted_at IS NULL AND o.status > 0").
		Group("product.category_id")

	err := db.Table("category").
		Where("deleted_at IS NULL").
		Joins("LEFT JOIN (?) AS subquery ON category.id = subquery.category_id", subquery).
		Select("category.id as category_id, category.name as category_name, subquery.price").
		Scan(&data).Error

	if err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	for _, item := range data {
		totalPrice += item.Price
	}

	for i, item := range data {
		if item.Price > 0 && totalPrice > 0 {
			item.Percentage = float32((item.Price / totalPrice) * 100)
			data[i] = item
		}
	}

	return ctx.JSON(data, http.StatusOK)
}
