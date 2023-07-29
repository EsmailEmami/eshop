package controllers

import (
	"net/http"

	"github.com/esmailemami/eshop/apphttp"
	appmodels "github.com/esmailemami/eshop/apphttp/models"
	"github.com/esmailemami/eshop/consts"
	"github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/errors"
	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
)

// GetCategorys godoc
// @Tags Categorys
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} []appmodels.CategoryOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /category [get]
func GetCategories(ctx *apphttp.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx).Model(&models.Category{})

	var data []appmodels.CategoryOutPutModel

	if err := baseDB.Find(&data).Error; err != nil {
		return err
	}

	return ctx.JSON(data, http.StatusOK)
}

// Create Category godoc
// @Tags Categorys
// @Accept json
// @Produce json
// @Security Bearer
// @Param category   body  appmodels.CategoryReqModel  true  "Category model"
// @Success 200 {object} httpmodels.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /category  [post]
func CreateCategory(ctx *apphttp.HttpContext) error {
	var inputModel appmodels.CategoryReqModel

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
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	return ctx.QuickResponse(consts.Created, http.StatusOK)
}

// Edit Category godoc
// @Tags Categorys
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Param category   body  appmodels.CategoryReqModel  true  "Category model"
// @Success 200 {object} httpmodels.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /category/edit/{id}  [post]
func EditCategory(ctx *apphttp.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	var inputModel appmodels.CategoryReqModel

	err = ctx.BlindBind(&inputModel)
	if err != nil {
		return err
	}

	err = inputModel.ValidateUpdate()
	if err != nil {
		return err
	}

	baseDB := db.MustGormDBConn(ctx)

	var dbModel models.Category

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if db.Exists(baseDB, &models.Category{}, "code = ? and id != ?", inputModel.Code, id) {
		return errors.NewValidationError(consts.ExistedCode, nil)
	}

	inputModel.MergeWithDBData(&dbModel)
	if baseDB.Save(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.QuickResponse(consts.Updated, http.StatusOK)
}

// Delete Category godoc
// @Tags Categorys
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} httpmodels.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /category/delete/{id}  [post]
func DeleteCategory(ctx *apphttp.HttpContext) error {
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
