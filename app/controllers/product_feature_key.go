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

// GetProductFeatureKeys godoc
// @Tags ProductFeatureKeys
// @Accept json
// @Produce json
// @Security Bearer
// @Param categoryId  query  string  false  "Product Feature Category ID"
// @Param searchTerm  query  string  false  "Search Term"
// @Success 200 {object} []appmodels.ProductFeatureKeyOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /productFeatureKey [get]
func GetProductFeatureKeys(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx).Model(&models.ProductFeatureKey{})

	categoryID, ok := ctx.GetParam("categoryId")
	if ok {
		baseDB = baseDB.Where("product_feature_category_id = ?", categoryID)
	}

	searchTerm, ok := ctx.GetParam("searchTerm")
	if ok {
		searchTerm = "%" + searchTerm + "%"
		baseDB = baseDB.Where("name LIKE ?", searchTerm)
	}

	var data []appmodels.ProductFeatureKeyOutPutModel

	if err := baseDB.Find(&data).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.JSON(data, http.StatusOK)
}

// GetProductFeatureKey godoc
// @Tags ProductFeatureKeys
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} appmodels.ProductFeatureKeyOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /productFeatureKey/{id} [get]
func GetProductFeatureKey(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx).Model(&models.ProductFeatureKey{})

	var data appmodels.ProductFeatureKeyOutPutModel

	if err := baseDB.First(&data, "id", id).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	return ctx.JSON(data, http.StatusOK)
}

// Create ProductFeatureKey godoc
// @Tags ProductFeatureKeys
// @Accept json
// @Produce json
// @Security Bearer
// @Param ProductFeatureKey   body  appmodels.ProductFeatureKeyReqModel  true  "ProductFeatureKey model"
// @Success 200 {object} httpmodels.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /productFeatureKey  [post]
func CreateProductFeatureKey(ctx *app.HttpContext) error {
	var inputModel appmodels.ProductFeatureKeyReqModel

	err := ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx)

	err = inputModel.ValidateCreate(baseDB)
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	if err := baseDB.Create(inputModel.ToDBModel()).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.QuickResponse(consts.Created, http.StatusOK)
}

// Edit ProductFeatureKey godoc
// @Tags ProductFeatureKeys
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Param ProductFeatureKey   body  appmodels.ProductFeatureKeyReqModel  true  "ProductFeatureKey model"
// @Success 200 {object} httpmodels.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /productFeatureKey/edit/{id}  [post]
func EditProductFeatureKey(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	var inputModel appmodels.ProductFeatureKeyReqModel

	err = ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	baseDB := db.MustGormDBConn(ctx)
	err = inputModel.ValidateUpdate(baseDB)
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	var dbModel models.ProductFeatureKey

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if db.Exists(
		baseDB,
		&models.ProductFeatureKey{},
		"name = ? and id != ?",
		inputModel.Name,
		id,
	) {
		return errors.NewValidationError(consts.ExistedTitle, nil)
	}

	inputModel.MergeWithDBData(&dbModel)
	if baseDB.Save(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.QuickResponse(consts.Updated, http.StatusOK)
}

// Delete ProductFeatureKey godoc
// @Tags ProductFeatureKeys
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} httpmodels.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /productFeatureKey/delete/{id}  [post]
func DeleteProductFeatureKey(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))
	if err != nil {
		return err
	}

	baseDB := db.MustGormDBConn(ctx).Model(&models.ProductFeatureKey{})

	var dbModel models.ProductFeatureKey

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if baseDB.Delete(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, nil)
	}

	return ctx.QuickResponse(consts.Deleted, http.StatusOK)
}
