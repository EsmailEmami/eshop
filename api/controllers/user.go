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
	dbmodels "github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GetUsers godoc
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param searchTerm  query  string  false  "search for item"
// @Success 200 {object} parameter.ListResponse[appmodels.UserOutPutModel]
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/user [get]
func GetUsers(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx).Model(&models.User{})

	parameter := parameter.New[appmodels.UserOutPutModel](ctx, baseDB)

	data, err := parameter.SearchColumns("first_name", "last_name", "mobile", "email", "username").
		SortDescending("created_at").
		Execute(baseDB)

	if err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.JSON(*data, http.StatusOK)
}

// GetUser godoc
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} appmodels.UserOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/user/{id} [get]
func GetUser(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx).Model(&models.User{})

	var data appmodels.UserOutPutModel

	if err := baseDB.First(&data, "id", id).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	return ctx.JSON(data, http.StatusOK)
}

// Create User godoc
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param User   body  appmodels.UserReqModel  true  "User model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/user  [post]
func CreateUser(ctx *app.HttpContext) error {
	var inputModel appmodels.UserReqModel

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

// Edit User godoc
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Param User   body  appmodels.UserReqModel  true  "User model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/user/edit/{id}  [post]
func EditUser(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	var inputModel appmodels.UserReqModel

	err = ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	err = inputModel.ValidateUpdate()
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	baseDB := db.MustGormDBConn(ctx)

	var dbModel models.User

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if db.Exists(baseDB, &models.User{}, "username = ? and id != ?", inputModel.Username, id) {
		return errors.NewValidationError(consts.UsernameAlreadyExists, nil)
	}

	inputModel.MergeWithDBData(&dbModel)
	if baseDB.Save(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.QuickResponse(consts.Updated, http.StatusOK)
}

// Delete User godoc
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/user/delete/{id}  [post]
func DeleteUser(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))
	if err != nil {
		return err
	}

	baseDB := db.MustGormDBConn(ctx).Model(&models.User{})

	var dbModel models.User

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if baseDB.Delete(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, nil)
	}

	return ctx.QuickResponse(consts.Deleted, http.StatusOK)
}

// Recovery Password godoc
// @Summary recovery user password.
// @Description recovery user password.
// @Tags Users
// @Accept json
// @Produce json
// @Param id  path  string  true  "Record ID"
// @Param RecoveryPasswordInput  body  models.RecoveryPasswordModel  true  "Recovery password model"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/user/recoveryPasword/{id} [post]
func AdminUserRecoveryPassword(ctx *app.HttpContext) error {
	userID, err := uuid.Parse(ctx.GetPathParam("id"))
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	db := db.MustGormDBConn(ctx)

	var input appmodels.RecoveryPasswordModel
	if err := ctx.BlindBind(&input); err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	if err := input.Validate(); err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	// encrypt password
	pass, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err := db.Model(&dbmodels.User{}).Where("id=?", userID).UpdateColumn("password", string(pass)).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.QuickResponse(consts.OperationDone, http.StatusOK)
}
