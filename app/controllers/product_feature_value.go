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

// GetProductFeatureValues godoc
// @Tags ProductFeatureValues
// @Accept json
// @Produce json
// @Security Bearer
// @Param keyId  query  string  false  "Product Feature Key ID"
// @Param productId  query  string  false  "Product ID"
// @Param searchTerm  query  string  false  "Search Term"
// @Success 200 {object} []appmodels.ProductFeatureValueOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /productFeatureValue [get]
func GetProductFeatureValues(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx).Model(&models.ProductFeatureValue{})

	categoryID, ok := ctx.GetParam("keyId")
	if ok {
		baseDB = baseDB.Where("product_feature_key_id = ?", categoryID)
	}

	productID, ok := ctx.GetParam("productId")
	if ok {
		baseDB = baseDB.Where("product_id = ?", productID)
	}

	searchTerm, ok := ctx.GetParam("searchTerm")
	if ok {
		searchTerm = "%" + searchTerm + "%"
		baseDB = baseDB.Where("value LIKE ?", searchTerm)
	}

	var data []appmodels.ProductFeatureValueOutPutModel

	if err := baseDB.Find(&data).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.JSON(data, http.StatusOK)
}

// GetProductFeatureValue godoc
// @Tags ProductFeatureValues
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} appmodels.ProductFeatureValueOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /productFeatureValue/{id} [get]
func GetProductFeatureValue(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx).Model(&models.ProductFeatureValue{})

	var data appmodels.ProductFeatureValueOutPutModel

	if err := baseDB.First(&data, "id", id).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	return ctx.JSON(data, http.StatusOK)
}

// Create ProductFeatureValue godoc
// @Tags ProductFeatureValues
// @Accept json
// @Produce json
// @Security Bearer
// @Param productId path  string  true  "Product ID"
// @Param ProductFeatureValue   body  []appmodels.ProductFeatureValueReqModel  true  "ProductFeatureValue model"
// @Success 200 {object} httpmodels.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /productFeatureValue/{productId}  [post]
func CreateProductFeatureValue(ctx *app.HttpContext) error {
	productId, err := uuid.Parse(ctx.GetPathParam("productId"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	var inputModels []appmodels.ProductFeatureValueReqModel

	err = ctx.BlindBind(&inputModels)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx)
	baseTx := baseDB.Begin()

	// remove all last items if exists
	if err := baseTx.Where("product_id = ?", productId).Delete(&models.ProductFeatureValue{}).Error; err != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	for _, inputModel := range inputModels {
		inputModel.ProductID = productId
		err = inputModel.ValidateCreate(baseDB)
		if err != nil {
			baseTx.Rollback()
			return errors.NewValidationError(consts.ValidationError, err)
		}

		if err := baseTx.Create(inputModel.ToDBModel()).Error; err != nil {
			baseTx.Rollback()
			return errors.NewInternalServerError(consts.InternalServerError, err)
		}
	}

	baseTx.Commit()

	return ctx.QuickResponse(consts.Created, http.StatusOK)
}

// Delete ProductFeatureValue godoc
// @Tags ProductFeatureValues
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} httpmodels.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /productFeatureValue/delete/{id}  [post]
func DeleteProductFeatureValue(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))
	if err != nil {
		return err
	}

	baseDB := db.MustGormDBConn(ctx).Model(&models.ProductFeatureValue{})

	var dbModel models.ProductFeatureValue

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if baseDB.Delete(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, nil)
	}

	return ctx.QuickResponse(consts.Deleted, http.StatusOK)
}
