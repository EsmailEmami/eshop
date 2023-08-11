package controllers

import (
	"net/http"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/errors"
	appmodels "github.com/esmailemami/eshop/app/models"
	"github.com/esmailemami/eshop/app/services/order"
	"github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
)

// Create OrderItem godoc
// @Tags OrderItems
// @Accept json
// @Produce json
// @Security Bearer
// @Param OrderItem   body  appmodels.OrderItemReqModel  true  "OrderItem model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/orderItem  [post]
func CreateOrderItem(ctx *app.HttpContext) error {
	var inputModel appmodels.OrderItemReqModel

	err := ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	user, err := ctx.GetUser()
	if err != nil {
		return errors.NewUnauthorizedError(err.Error(), err)
	}

	baseDB := db.MustGormDBConn(ctx)
	baseTx := baseDB.Begin()

	err = inputModel.ValidateCreate(baseDB)
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	// get order
	order, err := order.GetOpenOrder(baseDB, *user.ID)
	if err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	// get order item quantity if exists
	var lastQuantity int
	baseDB.Model(&models.OrderItem{}).Select("quantity").First(&lastQuantity, "order_id = ? AND product_item_id = ?", *order.ID, inputModel.ProductItemID).Delete(&models.OrderItem{})

	// get product item
	productItem := struct {
		Price    float64 `gorm:"column:price"`
		Quantity int     `gorm:"column:quantity"`
	}{}

	if err := baseDB.Model(&models.ProductItem{}).
		Select("price, quantity").
		First(&productItem, "id=?", inputModel.ProductItemID).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	productItem.Quantity = productItem.Quantity + lastQuantity - inputModel.Quantity

	if productItem.Quantity < 0 {
		return errors.NewBadRequestError(consts.InvalidQuantity, nil)
	}

	// check if order item is existed delete it at first
	if err := baseTx.Unscoped().Model(&models.OrderItem{}).Where("order_id = ? AND product_item_id = ?", *order.ID, inputModel.ProductItemID).Delete(&models.OrderItem{}).Error; err != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	inputModel.OrderID = *order.ID
	inputModel.Price = productItem.Price
	dbModel := inputModel.ToDBModel()
	if err := baseTx.Create(dbModel).Error; err != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	// update productItem quantity
	baseTx.Model(&models.ProductItem{}).Where("id = ?", inputModel.ProductItemID).
		UpdateColumn("quantity", productItem.Quantity)

	baseTx.Commit()

	return ctx.QuickResponse(consts.Created, http.StatusOK)
}

// Delete OrderItem godoc
// @Tags OrderItems
// @Accept json
// @Produce json
// @Security Bearer
// @Param productItemId  path  string  true  "Product Item ID"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/orderItem/delete/{productItemId}  [post]
func DeleteOrderItem(ctx *app.HttpContext) error {
	productItemID, err := uuid.Parse(ctx.GetPathParam("productItemId"))
	if err != nil {
		return err
	}

	baseDB := db.MustGormDBConn(ctx)
	baseTx := baseDB.Begin()

	// get order item
	var dbModel models.OrderItem
	if baseDB.Preload("Order").Model(&models.OrderItem{}).Order("created_at DESC").First(&dbModel, "product_item_id", productItemID).Error != nil {
		baseTx.Rollback()
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if dbModel.Order.Status != models.OrderStatusOpen {
		return errors.NewBadRequestError("You do not have permission to delete the product.", nil)
	}

	var totalItems int64

	if err := baseDB.Model(&models.OrderItem{}).Where("order_id = ? and id != ?", dbModel.OrderID, *dbModel.ID).Count(&totalItems).Error; err != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	if baseTx.Unscoped().Delete(&dbModel).Error != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, nil)
	}

	// delete order if it is the last item
	if totalItems == 0 {
		if baseTx.Unscoped().Model(&models.Order{}).Where("id = ?", dbModel.OrderID).Delete(&models.Order{}).Error != nil {
			baseTx.Rollback()
			return errors.NewInternalServerError(consts.InternalServerError, nil)
		}
	}

	var productItemQuantity int
	if baseDB.Model(&models.ProductItem{}).Select("quantity").First(&productItemQuantity, dbModel.ProductItemID).Error != nil {
		baseTx.Rollback()
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	productItemQuantity = productItemQuantity + dbModel.Quantity

	// update productItem quantity
	baseTx.Model(&models.ProductItem{}).Where("id = ?", dbModel.ProductItemID).
		UpdateColumn("quantity", productItemQuantity)

	baseTx.Commit()

	return ctx.QuickResponse(consts.Deleted, http.StatusOK)
}
