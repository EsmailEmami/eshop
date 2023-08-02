package controllers

import (
	"net/http"

	"github.com/esmailemami/eshop/app"
	appmodels "github.com/esmailemami/eshop/app/models"
	"github.com/esmailemami/eshop/consts"
	"github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/errors"
	"github.com/esmailemami/eshop/models"
	"github.com/esmailemami/eshop/services/order"
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
// @Router /orderItem  [post]
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

	if productItem.Quantity-inputModel.Quantity < 0 {
		return errors.NewBadRequestError("تعداد وارد شده مجاز نمیباشد", nil)
	}

	inputModel.OrderID = *order.ID
	inputModel.Price = productItem.Price
	dbModel := inputModel.ToDBModel()
	if err := baseTx.Create(dbModel).Error; err != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	baseTx.Commit()

	return ctx.QuickResponse(consts.Created, http.StatusOK)
}
