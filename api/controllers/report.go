package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/errors"
	"github.com/esmailemami/eshop/app/helpers"
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

// GetSellsChart godoc
// @Tags Reports
// @Accept json
// @Produce json
// @Security Bearer
// @Param type query string true "chart type" Enums(hour,daily,weekly,monthly)
// @Param length  query  int  false  "rows length default is 15"
// @Success 200 {object} []helpers.KeyValueResponse[string, int]
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/report/sellsChart [get]
func ReportSellsChart(ctx *app.HttpContext) error {
	baseDB := dbpkg.MustGormDBConn(ctx)

	baseDB = baseDB.Table("order_item oi").
		Joins(`INNER JOIN "order" o ON o.id = oi.order_id`).
		Where("o.deleted_at IS NULL AND o.status > 0 AND o.paid_at IS NOT NULL")

	queryType, ok := ctx.GetParam("type")
	if !ok {
		queryType = "daily"
	}

	length, ok := ctx.GetParam("length")
	if !ok {
		length = "15"
	}

	lengthNum, err := strconv.Atoi(length)
	if err != nil {
		lengthNum = 15
	}

	_ = length
	switch queryType {
	case "hour":
		baseDB = baseDB.Where("o.paid_at >= ?", time.Now().Add(time.Duration(-lengthNum)*time.Hour)).
			Select(`TO_CHAR(DATE_TRUNC('hour', o.paid_at), 'HH24:MI') AS "key", SUM(oi.quantity) AS "value"`)
	case "daily":
		baseDB = baseDB.Where("o.paid_at >= ?", time.Now().AddDate(0, 0, -lengthNum)).
			Select(`TO_CHAR(DATE_TRUNC('day', o.paid_at), 'YYYY-MM-DD') AS "key", SUM(oi.quantity) AS "value"`)
	case "weekly":
		baseDB = baseDB.Where("o.paid_at >= ?", time.Now().AddDate(0, 0, -lengthNum*7)).
			Select(`TO_CHAR(DATE_TRUNC('week', o.paid_at), 'YYYY-MM-DD') AS "key", SUM(oi.quantity) AS "value"`)
	case "monthly":
		baseDB = baseDB.Where("o.paid_at >= ?", time.Now().AddDate(0, -lengthNum, 0)).
			Select(`TO_CHAR(DATE_TRUNC('month', o.paid_at), 'YYYY-MM-DD') AS "key", SUM(oi.quantity) AS "value"`)
	default:
		return errors.NewBadRequestError("Invalid type", nil)
	}

	var data []helpers.KeyValueResponse[string, int]

	if err := baseDB.Group(`"key"`).Order(`"key" ASC`).Find(&data).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.JSON(data, http.StatusOK)
}
