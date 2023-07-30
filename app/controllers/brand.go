package controllers

import (
	"net/http"

	"github.com/esmailemami/eshop/app"
	appmodels "github.com/esmailemami/eshop/app/models"
	"github.com/esmailemami/eshop/consts"
	"github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/errors"
	"github.com/esmailemami/eshop/models"
	fileService "github.com/esmailemami/eshop/services/file"
	"github.com/google/uuid"
)

// GetBrands godoc
// @Tags Brands
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} []appmodels.BrandOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /brand [get]
func GetBrands(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx).Model(&models.Brand{})

	var data []appmodels.BrandOutPutModel

	if err := baseDB.Find(&data).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.JSON(data, http.StatusOK)
}

// GetBrand godoc
// @Tags Brands
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} appmodels.BrandOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /brand/{id} [get]
func GetBrand(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx).Model(&models.Brand{})

	var data appmodels.BrandOutPutModel

	if err := baseDB.First(&data, "id", id).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	return ctx.JSON(data, http.StatusOK)
}

// Create Brand godoc
// @Tags Brands
// @Accept json
// @Produce json
// @Security Bearer
// @Param Brand   body  appmodels.BrandReqModel  true  "Brand model"
// @Success 200 {object} httpmodels.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /brand  [post]
func CreateBrand(ctx *app.HttpContext) error {
	var inputModel appmodels.BrandReqModel

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

// Edit Brand godoc
// @Tags Brands
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Param Brand   body  appmodels.BrandReqModel  true  "Brand model"
// @Success 200 {object} httpmodels.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /brand/edit/{id}  [post]
func EditBrand(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	var inputModel appmodels.BrandReqModel

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

	var dbModel models.Brand

	if baseDB.Preload("File").First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if db.Exists(baseDB, &models.Brand{}, "code = ? and id != ?", inputModel.Code, id) {
		return errors.NewValidationError(consts.ExistedCode, nil)
	}

	if dbModel.FileID != inputModel.FileID {
		err := fileService.DeleteFile(baseTx, dbModel.File)

		if err != nil {
			baseTx.Rollback()
			return errors.NewInternalServerError(consts.InternalServerError, err)
		}
	}

	inputModel.MergeWithDBData(&dbModel)
	if baseTx.Save(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	baseTx.Commit()

	return ctx.QuickResponse(consts.Updated, http.StatusOK)
}

// Delete Brand godoc
// @Tags Brands
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} httpmodels.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /brand/delete/{id}  [post]
func DeleteBrand(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))
	if err != nil {
		return err
	}

	baseDB := db.MustGormDBConn(ctx)
	baseTx := baseDB.Begin()

	var dbModel models.Brand

	if baseDB.Preload("File").First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	err = fileService.DeleteFile(baseTx, dbModel.File)

	if err != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	if baseTx.Delete(&dbModel).Error != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, nil)
	}

	baseTx.Commit()

	return ctx.QuickResponse(consts.Deleted, http.StatusOK)
}
