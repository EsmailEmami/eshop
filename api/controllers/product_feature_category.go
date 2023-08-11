package controllers

import (
	"net/http"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/errors"
	appmodels "github.com/esmailemami/eshop/app/models"
	"github.com/esmailemami/eshop/app/parameter"
	"github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
)

// GetProductFeatureCategories godoc
// @Tags ProductFeatureCategories
// @Accept json
// @Produce json
// @Security Bearer
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param searchTerm  query  string  false  "search for item"
// @Success 200 {object} parameter.ListResponse[appmodels.ProductFeatureCategoryOutPutModel]
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/productFeatureCategory [get]
func GetProductFeatureCategories(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx).Model(&models.ProductFeatureCategory{})

	parameter := parameter.New[appmodels.ProductFeatureCategoryOutPutModel](ctx, baseDB)

	data, err := parameter.SearchColumns("name", "code").
		SortDescending("created_at").
		Execute(baseDB)

	if err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.JSON(*data, http.StatusOK)
}

// GetProductFeatureCategory godoc
// @Tags ProductFeatureCategories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} appmodels.ProductFeatureCategoryOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/productFeatureCategory/{id} [get]
func GetProductFeatureCategory(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx).Model(&models.ProductFeatureCategory{})

	var data appmodels.ProductFeatureCategoryOutPutModel

	if err := baseDB.First(&data, "id", id).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	return ctx.JSON(data, http.StatusOK)
}

// Create ProductFeatureCategory godoc
// @Tags ProductFeatureCategories
// @Accept json
// @Produce json
// @Security Bearer
// @Param ProductFeatureCategory   body  appmodels.ProductFeatureCategoryReqModel  true  "ProductFeatureCategory model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/productFeatureCategory  [post]
func CreateProductFeatureCategory(ctx *app.HttpContext) error {
	var inputModel appmodels.ProductFeatureCategoryReqModel

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

// Edit ProductFeatureCategory godoc
// @Tags ProductFeatureCategories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Param ProductFeatureCategory   body  appmodels.ProductFeatureCategoryReqModel  true  "ProductFeatureCategory model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/productFeatureCategory/edit/{id}  [post]
func EditProductFeatureCategory(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	var inputModel appmodels.ProductFeatureCategoryReqModel

	err = ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	err = inputModel.ValidateUpdate()
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	baseDB := db.MustGormDBConn(ctx)

	var dbModel models.ProductFeatureCategory

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if db.Exists(
		baseDB,
		&models.ProductFeatureCategory{},
		"code = ? and id != ?",
		inputModel.Code,
		id,
	) {
		return errors.NewValidationError(consts.ExistedCode, nil)
	}

	if db.Exists(
		baseDB,
		&models.ProductFeatureCategory{},
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

// Delete ProductFeatureCategory godoc
// @Tags ProductFeatureCategories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/productFeatureCategory/delete/{id}  [post]
func DeleteProductFeatureCategory(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))
	if err != nil {
		return err
	}

	baseDB := db.MustGormDBConn(ctx).Model(&models.ProductFeatureCategory{})

	var dbModel models.ProductFeatureCategory

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if baseDB.Delete(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, nil)
	}

	return ctx.QuickResponse(consts.Deleted, http.StatusOK)
}
