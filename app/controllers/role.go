package controllers

import (
	"net/http"

	"github.com/esmailemami/eshop/app"
	appmodels "github.com/esmailemami/eshop/app/models"
	"github.com/esmailemami/eshop/app/parameter"
	"github.com/esmailemami/eshop/consts"
	"github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/errors"
	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
)

// GetRoles godoc
// @Tags Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param searchTerm  query  string  false  "search for item"
// @Success 200 {object} parameter.ListResponse[appmodels.RoleOutPutModel]
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/role [get]
func GetRoles(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx).Model(&models.Role{})

	parameter := parameter.New[appmodels.RoleOutPutModel](ctx, baseDB)

	data, err := parameter.SearchColumns("name", "code").
		SortDescending("created_at").
		Execute(baseDB)

	if err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.JSON(*data, http.StatusOK)
}

// GetRole godoc
// @Tags Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} appmodels.RoleOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/role/{id} [get]
func GetRole(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx).Model(&models.Role{})

	var data appmodels.RoleOutPutModel

	if err := baseDB.First(&data, "id", id).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	return ctx.JSON(data, http.StatusOK)
}

// Create Role godoc
// @Tags Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param Role   body  appmodels.RoleReqModel  true  "Role model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/role  [post]
func CreateRole(ctx *app.HttpContext) error {
	var inputModel appmodels.RoleReqModel

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

// Edit Role godoc
// @Tags Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Param Role   body  appmodels.RoleReqModel  true  "Role model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/role/edit/{id}  [post]
func EditRole(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	var inputModel appmodels.RoleReqModel

	err = ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	err = inputModel.ValidateUpdate()
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	baseDB := db.MustGormDBConn(ctx)

	var dbModel models.Role

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if db.Exists(baseDB, &models.Role{}, "code = ? and id != ?", inputModel.Code, id) {
		return errors.NewValidationError(consts.ExistedCode, nil)
	}

	inputModel.MergeWithDBData(&dbModel)
	if baseDB.Save(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.QuickResponse(consts.Updated, http.StatusOK)
}

// Delete Role godoc
// @Tags Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/role/delete/{id}  [post]
func DeleteRole(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))
	if err != nil {
		return err
	}

	baseDB := db.MustGormDBConn(ctx).Model(&models.Role{})

	var dbModel models.Role

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if baseDB.Delete(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, nil)
	}

	return ctx.QuickResponse(consts.Deleted, http.StatusOK)
}

// GetRole godoc
// @Tags Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} []models.Action
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/role/Permissions [get]
func GetPermissions(ctx *app.HttpContext) error {
	actions := models.GetPermissionsTree()
	return ctx.JSON(actions, http.StatusOK)
}
