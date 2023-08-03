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

// GetAddresses godoc
// @Tags Addresses
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} []appmodels.AddressOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /address [get]
func GetAddresses(ctx *app.HttpContext) error {
	user, err := ctx.GetUser()
	if err != nil {
		return errors.NewUnauthorizedError(consts.UnauthorizedError, err)
	}

	baseDB := db.MustGormDBConn(ctx).Model(&models.Address{}).Where("created_by_id=?", *user.ID)

	var data []appmodels.AddressOutPutModel

	if err := baseDB.Find(&data).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.JSON(data, http.StatusOK)
}

// GetAddress godoc
// @Tags Addresses
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} appmodels.AddressOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /address/{id} [get]
func GetAddress(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx).Model(&models.Address{})

	var data appmodels.AddressOutPutModel

	if err := baseDB.First(&data, "id", id).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	return ctx.JSON(data, http.StatusOK)
}

// Create Address godoc
// @Tags Addresses
// @Accept json
// @Produce json
// @Security Bearer
// @Param Address   body  appmodels.AddressReqModel  true  "Address model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /address  [post]
func CreateAddress(ctx *app.HttpContext) error {
	var inputModel appmodels.AddressReqModel

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

// Edit Address godoc
// @Tags Addresses
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Param Address   body  appmodels.AddressReqModel  true  "Address model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /address/edit/{id}  [post]
func EditAddress(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	var inputModel appmodels.AddressReqModel

	err = ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	err = inputModel.ValidateUpdate()
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	baseDB := db.MustGormDBConn(ctx)

	var dbModel models.Address

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	inputModel.MergeWithDBData(&dbModel)
	if baseDB.Save(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.QuickResponse(consts.Updated, http.StatusOK)
}

// Delete Address godoc
// @Tags Addresses
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /address/delete/{id}  [post]
func DeleteAddress(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))
	if err != nil {
		return err
	}

	baseDB := db.MustGormDBConn(ctx).Model(&models.Address{})

	var dbModel models.Address

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if baseDB.Delete(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, nil)
	}

	return ctx.QuickResponse(consts.Deleted, http.StatusOK)
}
