package controllers

import (
	"net/http"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/errors"
	appmodels "github.com/esmailemami/eshop/app/models"
	fileService "github.com/esmailemami/eshop/app/services/file"
	"github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
)

// GetAppPics godoc
// @Tags AppPics
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} []appmodels.AppPicOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/appPic [get]
func GetAppPics(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx)

	baseDB = baseDB.Table("app_pic ap").
		Joins("INNER JOIN file f ON f.id = ap.file_id").
		Where("ap.deleted_at IS NULL")

	var data []appmodels.AppPicOutPutModel

	if err := baseDB.Select("ap.id, ap.created_at, ap.updated_at, ap.priority, ap.app_pic_type,ap.file_id,ap.title,ap.description,ap.url,f.file_type,f.unique_file_name as file_name").Find(&data).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	for k, v := range data {
		v.FileUrl = v.FileType.GetFileUrl(v.FileName)
		data[k] = v
	}

	return ctx.JSON(data, http.StatusOK)
}

// GetAppPic godoc
// @Tags AppPics
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} appmodels.AppPicOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/appPic/{id} [get]
func GetAppPic(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx)

	baseDB = baseDB.Table("app_pic ap").
		Joins("INNER JOIN file f ON f.id = ap.file_id").
		Where("ap.deleted_at IS NULL")

	var data appmodels.AppPicOutPutModel

	if err := baseDB.Select("ap.id, ap.created_at, ap.updated_at, ap.priority, ap.app_pic_type,ap.file_id,ap.title,ap.description,ap.url,f.file_type,f.unique_file_name as file_name").
		First(&data, "ap.id=?", id).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	data.FileUrl = data.FileType.GetFileUrl(data.FileName)

	return ctx.JSON(data, http.StatusOK)
}

// Create AppPic godoc
// @Tags AppPics
// @Accept json
// @Produce json
// @Security Bearer
// @Param AppPic   body  appmodels.AppPicReqModel  true  "AppPic model"
// @Success 200 {object} helpers.SuccessDBResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/appPic  [post]
func CreateAppPic(ctx *app.HttpContext) error {
	var inputModel appmodels.AppPicReqModel

	err := ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx)

	err = inputModel.ValidateCreate()
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	dbModel := inputModel.ToDBModel()
	if err := baseDB.Create(dbModel).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.QuickDBResponse(consts.Created, *dbModel.ID, http.StatusOK)
}

// Edit AppPic godoc
// @Tags AppPics
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Param AppPic   body  appmodels.AppPicReqModel  true  "AppPic model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/appPic/edit/{id}  [post]
func EditAppPic(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	var inputModel appmodels.AppPicReqModel

	err = ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	baseDB := db.MustGormDBConn(ctx)
	baseTx := baseDB.Begin()

	var dbModel models.AppPic

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	err = inputModel.ValidateUpdate(id)
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	inputModel.MergeWithDBData(&dbModel)
	if baseTx.Save(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	baseTx.Commit()

	return ctx.QuickResponse(consts.Updated, http.StatusOK)
}

// Delete AppPic godoc
// @Tags AppPics
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/appPic/delete/{id}  [post]
func DeleteAppPic(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))
	if err != nil {
		return err
	}

	baseDB := db.MustGormDBConn(ctx)
	baseTx := baseDB.Begin()

	var dbModel models.AppPic

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
