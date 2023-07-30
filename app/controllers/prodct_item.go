package controllers

import (
	"net/http"

	"github.com/esmailemami/eshop/app"
	appmodels "github.com/esmailemami/eshop/app/models"
	"github.com/esmailemami/eshop/consts"
	"github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/errors"
	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
)

// GetProductItem godoc
// @Tags ProductItems
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} appmodels.ProductItemInfoOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /productItem/{id} [get]
func GetProductItem(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx).Model(&models.ProductItem{})

	var data appmodels.ProductItemInfoOutPutModel

	if err := baseDB.Table("product_item pi2").
		Joins("INNER JOIN product p ON P.id = pi2 .product_id").
		Joins("INNER JOIN color c ON C.id = pi2.color_id").
		Select(`pi2.id, pi2.price,pi2.status, pi2 .color_id, pi2.product_id, pi2.quantity,
		p."name" AS product_title, p.code AS product_code, p.short_description AS product_short_description, 
		p.description AS product_description, c."name" AS color_name`).
		First(&data, "pi2.id", id).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	return ctx.JSON(data, http.StatusOK)
}

// Create ProductItem godoc
// @Tags ProductItems
// @Accept json
// @Produce json
// @Security Bearer
// @Param ProductItem   body  appmodels.ProductItemReqModel  true  "ProductItem model"
// @Success 200 {object} httpmodels.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /productItem  [post]
func CreateProductItem(ctx *app.HttpContext) error {
	var inputModel appmodels.ProductItemReqModel

	err := ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx)
	baseTx := baseDB.Begin()

	err = inputModel.ValidateCreate(baseDB)
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	dbModel := inputModel.ToDBModel()
	if err := baseTx.Create(dbModel).Error; err != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	if inputModel.IsMainItem {
		if err := baseTx.Model(&models.Product{}).
			Where("id=?", inputModel.ProductID).
			UpdateColumn("default_product_item_id", dbModel.ID).Error; err != nil {
			baseTx.Rollback()
			return errors.NewInternalServerError(consts.InternalServerError, err)
		}
	}

	baseTx.Commit()

	return ctx.QuickResponse(consts.Created, http.StatusOK)
}

// Edit ProductItem godoc
// @Tags ProductItems
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Param ProductItem   body  appmodels.ProductItemReqModel  true  "ProductItem model"
// @Success 200 {object} httpmodels.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /productItem/edit/{id}  [post]
func EditProductItem(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	var inputModel appmodels.ProductItemReqModel

	err = ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	baseDB := db.MustGormDBConn(ctx)
	baseTx := baseDB.Begin()

	err = inputModel.ValidateUpdate(baseDB)
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	var dbModel models.ProductItem

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	inputModel.MergeWithDBData(&dbModel)
	if baseTx.Save(&dbModel).Error != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	if inputModel.IsMainItem {
		if err := baseTx.Model(&models.Product{}).
			Where("id=?", inputModel.ProductID).
			UpdateColumn("default_product_item_id", dbModel.ID).Error; err != nil {
			baseTx.Rollback()
			return errors.NewInternalServerError(consts.InternalServerError, err)
		}
	}

	baseTx.Commit()

	return ctx.QuickResponse(consts.Updated, http.StatusOK)
}

// Delete ProductItem godoc
// @Tags ProductItems
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} httpmodels.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /productItem/delete/{id}  [post]
func DeleteProductItem(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))
	if err != nil {
		return err
	}

	baseDB := db.MustGormDBConn(ctx)
	baseTx := baseDB.Begin()

	var dbModel models.ProductItem

	if baseDB.Preload("Product").First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if baseTx.Delete(&dbModel).Error != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, nil)
	}

	// choose last most sell item
	if dbModel.Product.DefaultProductItemID.String() == dbModel.ID.String() {
		mostSellDBItem := baseDB.Model(&models.ProductItem{}).
			Where("id != ?", dbModel.ID).Order("bought_quantity DESC").
			Select("TOP(1) id")

		if err := baseTx.Model(&models.Product{}).
			Where("id=?", dbModel.ProductID).
			UpdateColumn("default_product_item_id", mostSellDBItem).Error; err != nil {
			baseTx.Rollback()
			return errors.NewInternalServerError(consts.InternalServerError, err)
		}
	}

	baseTx.Commit()

	return ctx.QuickResponse(consts.Deleted, http.StatusOK)
}
