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

// Create FavoriteProductItem godoc
// @Tags FavoriteProductItems
// @Accept json
// @Produce json
// @Security Bearer
// @Param FavoriteProductItem   body  appmodels.FavoriteProductItemReqModel  true  "FavoriteProductItem model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /favoriteProductItem  [post]
func CreateFavoriteProductItem(ctx *app.HttpContext) error {
	user, err := ctx.GetUser()
	if err != nil {
		return errors.NewUnauthorizedError(consts.UnauthorizedError, err)
	}

	var inputModel appmodels.FavoriteProductItemReqModel

	err = ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx)

	err = inputModel.ValidateCreate(baseDB)
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	if db.Exists(baseDB, &models.FavoriteProductItem{}, "created_by_id=? AND product_item_id=?", *user.ID, inputModel.ProductItemID) {
		return errors.NewBadRequestError("The selected product already exists in your favorites", nil)
	}

	if err := baseDB.Create(inputModel.ToDBModel()).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.QuickResponse(consts.Created, http.StatusOK)
}

// Delete FavoriteProductItem godoc
// @Tags FavoriteProductItems
// @Accept json
// @Produce json
// @Security Bearer
// @Param productItemId  path  string  true  "Product Item ID"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /favoriteProductItem/delete/{productItemId}  [post]
func DeleteFavoriteProductItem(ctx *app.HttpContext) error {
	user, err := ctx.GetUser()
	if err != nil {
		return errors.NewUnauthorizedError(consts.UnauthorizedError, err)
	}

	productItemID, err := uuid.Parse(ctx.GetPathParam("productItemId"))
	if err != nil {
		return err
	}

	baseDB := db.MustGormDBConn(ctx)
	baseTx := baseDB.Begin()

	var dbModel models.FavoriteProductItem

	if baseDB.First(&dbModel, "created_by_id=? AND product_item_id=?", *user.ID, productItemID).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if baseTx.Where("id=?", *dbModel.ID).Delete(&models.FavoriteProductItem{}).Error != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, nil)
	}

	baseTx.Commit()

	return ctx.QuickResponse(consts.Deleted, http.StatusOK)
}
