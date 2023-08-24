package controllers

import (
	"net/http"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/errors"
	"github.com/esmailemami/eshop/app/helpers"
	appmodels "github.com/esmailemami/eshop/app/models"
	"github.com/esmailemami/eshop/app/parameter"
	"github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
)

// GetCategories godoc
// @Tags Categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param searchTerm  query  string  false  "search for item"
// @Success 200 {object} parameter.ListResponse[appmodels.CategoryOutPutModel]
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/category [get]
func GetCategories(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx).Model(&models.Category{})

	parameter := parameter.New[appmodels.CategoryOutPutModel](ctx, baseDB)

	data, err := parameter.SearchColumns("name", "code").
		SortDescending("created_at").
		Execute(baseDB)

	if err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.JSON(*data, http.StatusOK)
}

// GetCategoriesSelectList godoc
// @Tags Categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param searchTerm  query  string  false  "search for item"
// @Success 200 {object} parameter.ListResponse[helpers.KeyValueResponse[uuid.UUID, string]]
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/category/selectList [get]
func GetCategoriesSelectList(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx)

	parameter := parameter.New[helpers.KeyValueResponse[uuid.UUID, string]](ctx, baseDB)

	baseDB = baseDB.Model(&models.Category{})

	response, err := parameter.SelectColumns(`id as "key", "name" as "value"`).
		SearchColumns(`name`).
		SortDescending("created_at", `name`).
		Execute(baseDB)

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	return ctx.JSON(*response, http.StatusOK)
}

// GetCategory godoc
// @Tags Categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} appmodels.CategoryOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/category/{id} [get]
func GetCategory(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx).Model(&models.Category{})

	var data appmodels.CategoryOutPutModel

	if err := baseDB.First(&data, "id", id).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	return ctx.JSON(data, http.StatusOK)
}

// Create Category godoc
// @Tags Categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param category   body  appmodels.CategoryReqModel  true  "Category model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/category  [post]
func CreateCategory(ctx *app.HttpContext) error {
	var inputModel appmodels.CategoryReqModel

	err := ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx)

	err = inputModel.ValidateCreate()
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	if err := baseDB.Create(inputModel.ToDBModel()).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.QuickResponse(consts.Created, http.StatusOK)
}

// Edit Category godoc
// @Tags Categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Param category   body  appmodels.CategoryReqModel  true  "Category model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/category/edit/{id}  [post]
func EditCategory(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	var inputModel appmodels.CategoryReqModel

	err = ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	baseDB := db.MustGormDBConn(ctx)

	var dbModel models.Category

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	err = inputModel.ValidateUpdate(id)
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	inputModel.MergeWithDBData(&dbModel)
	if baseDB.Save(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.QuickResponse(consts.Updated, http.StatusOK)
}

// Delete Category godoc
// @Tags Categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/category/delete/{id}  [post]
func DeleteCategory(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))
	if err != nil {
		return err
	}

	baseDB := db.MustGormDBConn(ctx).Model(&models.Category{})

	var dbModel models.Category

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if baseDB.Delete(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, nil)
	}

	return ctx.QuickResponse(consts.Deleted, http.StatusOK)
}
