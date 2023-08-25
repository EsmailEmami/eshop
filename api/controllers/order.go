package controllers

import (
	"net/http"
	"time"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/errors"
	appmodels "github.com/esmailemami/eshop/app/models"
	"github.com/esmailemami/eshop/app/parameter"
	"github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
)

// GetOrders godoc
// @Tags Orders
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} appmodels.OrderOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/order [get]
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
		Joins("LEFT JOIN LATERAL (?) as d ON TRUE", baseDB.Table("discount d").
			Where("d.product_item_id = pi2.id AND d.deleted_at IS NULL AND d.type = 1").
			Where("CASE WHEN d.expires_in IS NOT NULL THEN d.expires_in > NOW() WHEN d.quantity IS NOT NULL THEN d.quantity > 0 ELSE TRUE END").
			Where("d.related_user_id IS NULL").
			Order("d.created_at ASC").
			Select("d.type, d.value, d.product_item_id, d.quantity").
			Limit(1),
		).
		Joins("INNER JOIN product p ON	p.id = pi2.product_id").
		Joins("CROSS JOIN LATERAL (?) as pf", baseDB.Table("product_file_map pfm").
			Select("file_id").
			Where("pfm.product_id = p.id").
			Order("pfm.priority ASC").
			Limit(1),
		).
		Joins("INNER JOIN file f ON f.id = pf.file_id").
		Where("o.status = 0 AND o.created_by_id = ?", *user.ID).
		Select(`oi.id, oi.product_item_id,p."name" AS product_name,pi2.price,oi.quantity,f.file_type, f.unique_file_name AS file_name,d.type as discount_type, d.value as discount_value, d.quantity as discount_quantity`).
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

// CheckoutOrder godoc
// @Tags Orders
// @Accept json
// @Produce json
// @Security Bearer
// @Param addressId  path  string  true  "Record ID"
// @Param discountId  query  string  false  "Discount ID"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/order/checkout/{addressId} [post]
func CheckoutOrder(ctx *app.HttpContext) error {
	user, err := ctx.GetUser()
	if err != nil {
		return errors.NewUnauthorizedError(err.Error(), err)
	}
	addressID, err := uuid.Parse(ctx.GetPathParam("addressId"))
	if err != nil {
		return errors.NewBadRequestError(consts.ModelAddressNotFound, err)
	}

	var (
		baseDB                = db.MustGormDBConn(ctx)
		baseTx                = baseDB.Begin()
		discountValue float64 = 0
		discountPrice float64 = 0
		discountType  models.DiscountType
	)

	if discountID, ok := ctx.GetParam("discountId"); ok {
		var dbModel models.Discount

		if baseDB.First(&dbModel, discountID).Error != nil {
			return errors.NewRecordNotFoundError(consts.ModelDiscountNotFound, nil)
		}

		if err := dbModel.IsValidToUse(user.ID, nil); err != nil {
			baseTx.Rollback()
			return errors.NewBadRequestError(err.Error(), err)
		}

		if dbModel.Quantity != nil {
			qty := *dbModel.Quantity - 1
			dbModel.Quantity = &qty

			if err := baseDB.Save(&dbModel).Error; err != nil {
				baseTx.Rollback()
				return errors.NewInternalServerError(consts.InternalServerError, err)
			}
		}

		discountValue = dbModel.Value
		discountType = dbModel.Type
	}

	var order models.Order

	if err := baseDB.Model(&models.Order{}).First(&order, "status = 0 AND created_by_id = ?", *user.ID).Error; err != nil {
		baseTx.Rollback()
		return errors.NewRecordNotFoundError(consts.ModelOrderNotFound, err)
	}

	var address models.Address

	if err := baseDB.Model(&models.Address{}).First(&address, "id=? AND created_by_id=?", addressID, *user.ID).Error; err != nil {
		baseTx.Rollback()
		return errors.NewRecordNotFoundError(consts.ModelAddressNotFound, err)
	}

	// update the prices of order items
	baseTx.Table("order_item oi").Where("oi.order_id=?", *order.ID).UpdateColumns(map[string]interface{}{
		"total_price": baseDB.Model(&models.ProductItem{}).
			Select("price").
			Where("id = oi.product_item_id").
			Limit(1),
		"price": baseDB.Table("product_item pi2").
			Select("pi2.price - (CASE WHEN d.type = 1 THEN ((d.value / 100) * pi2.price) WHEN d.type = 0 THEN d.value  ELSE 0 END)").
			Where("id = oi.product_item_id").
			Joins("LEFT JOIN LATERAL (?) as d ON TRUE", baseDB.Table("discount d").
				Where("d.product_item_id = pi2.id AND d.deleted_at IS NULL AND d.type = 1").
				Where("CASE WHEN d.expires_in IS NOT NULL THEN d.expires_in > NOW() WHEN d.quantity IS NOT NULL THEN d.quantity > 0 ELSE TRUE END").
				Where("d.related_user_id IS NULL").
				Order("d.created_at ASC").
				Select("d.type, d.value, d.product_item_id").
				Limit(1),
			).
			Limit(1),
	})

	orderItemsPrice := struct {
		Price      float64 `gorm:"price"`
		TotalPrice float64 `gorm:"total_price"`
	}{}

	if err := baseTx.Model(&models.OrderItem{}).Select("SUM(price) as price, SUM(total_price) as total_price").
		Where("order_id=?", *order.ID).Find(&orderItemsPrice).Error; err != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	if discountValue > 0 {
		switch discountType {
		case models.DiscountTypeNumeric:
			discountPrice = discountValue
			orderItemsPrice.Price -= discountPrice
		case models.DiscountTypePercent:
			discountPrice = orderItemsPrice.Price * (discountValue / 100)
			orderItemsPrice.Price -= discountPrice
		}
	}

	var (
		orderDiscountPrice *float64             = nil
		orderDiscountType  *models.DiscountType = nil
		orderDiscountValue *float64             = nil
	)

	if discountValue > 0 {
		orderDiscountPrice = &discountPrice
		orderDiscountType = &discountType
		orderDiscountValue = &discountValue
	}

	// update order
	if err := baseTx.Model(&models.Order{}).
		Where("id=?", *order.ID).UpdateColumns(map[string]interface{}{
		"status":         models.OrderStatusPaid,
		"price":          orderItemsPrice.Price,
		"discount_price": orderDiscountPrice,
		"discount_value": orderDiscountValue,
		"discount_type":  orderDiscountType,
		"total_price":    orderItemsPrice.TotalPrice,
		"first_name":     address.FirstName,
		"last_name":      address.LastName,
		"plaque":         address.Plaque,
		"phone_number":   address.PhoneNumber,
		"national_code":  address.NationalCode,
		"postal_code":    address.PostalCode,
		"address":        address.Address,
		"paid_at":        time.Now(),
	}).Error; err != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	if err := baseTx.Commit().Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.QuickResponse(consts.OperationDone, http.StatusOK)
}

// GetOrdersForAdmin godoc
// @Tags Orders
// @Accept json
// @Produce json
// @Security Bearer
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param searchTerm  query  string  false  "search for item"
// @Param status  query  int  false  "order status"
// @Success 200 {object} parameter.ListResponse[appmodels.AdminOrderOutPutModel]
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/order [get]
func GetAdminOrders(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx)
	parameter := parameter.New[appmodels.AdminOrderOutPutModel](ctx, baseDB)

	baseDB = baseDB.Table(`"order" o`).
		Joins(`INNER JOIN "user" u ON u.id = o.created_by_id`)

	if status, ok := ctx.GetParam("status"); ok {
		baseDB = baseDB.Where("o.status = ?", status)
	}

	data, err := parameter.SearchColumns("u.username", "u.first_name", "u.last_name").
		SortDescending("o.updated_at").
		SelectColumns(`u.id as user_id, o.id as order_id, u.username, o.status, o.created_at, o.updated_at, 
			CASE
				WHEN o.status > 0 THEN o.price
				ELSE (SELECT SUM(price) FROM order_item oi WHERE oi.order_id = o.id)
			END AS price,
			CASE
				WHEN o.status > 0 THEN o.total_price
				ELSE (SELECT SUM(total_price) FROM order_item oi WHERE oi.order_id = o.id)
			END AS total_price,
			o.paid_at, o.discount_price,o.discount_value,o.discount_type
		`).
		Execute(baseDB)

	if err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.JSON(data, http.StatusOK)
}
