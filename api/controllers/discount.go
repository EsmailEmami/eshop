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
	"gorm.io/gorm"
)

// GetDiscounts godoc
// @Tags Discounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param searchTerm  query  string  false  "search for item"
// @Param productItemId  query  string  false  "search for product item"
// @Param relatedUser  query  string  false  "search for related user"
// @Param createdBy  query  string  false  "search for creator"
// @Success 200 {object} parameter.ListResponse[appmodels.DiscountAdminOutPutModel]
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/discount [get]
func GetDiscounts(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx)
	parameter := parameter.New[appmodels.DiscountAdminOutPutModel](ctx, baseDB)

	baseDB = baseDB.Table("discount d").
		Joins("LEFT JOIN product_item pi2 ON pi2.id = d.product_item_id").
		Joins(`INNER JOIN "user" cu ON cu.id = d.created_by_id`).
		Joins(`LEFT JOIN "product" p ON p.id = pi2.product_id`).
		Joins(`LEFT JOIN "user" ru ON ru.id = d.related_user_id`).
		Where("d.deleted_at IS NULL")

	if productItemID, ok := ctx.GetParam("productItemId"); ok {
		baseDB = baseDB.Where("pi2.id = ?", productItemID)
	}

	if relatedUserID, ok := ctx.GetParam("relatedUser"); ok {
		baseDB = baseDB.Where("ru.id = ?", relatedUserID)
	}

	if creatorUserID, ok := ctx.GetParam("createdBy"); ok {
		baseDB = baseDB.Where("cu.id = ?", creatorUserID)
	}

	data, err := parameter.SelectColumns("d.id, d.created_at, d.updated_at, d.product_item_id, p.name as product_name, d.type, d.value, d.quantity, d.expires_in, d.code, d.created_by_id as creator_user_id, cu.username,d.related_user_id, ru.username as related_user_username").
		SortDescending("d.created_at", "d.expires_in").
		SearchColumns("ru.username, ru.first_name, ru.last_name", "cu.username, cu.first_name, cu.last_name", "d.code", "p.name", "p.code").
		EachItemProcess(func(db *gorm.DB, item *appmodels.DiscountAdminOutPutModel) error {
			item.TypeName = item.Type.String()
			return nil
		}).
		Execute(baseDB)

	if err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.JSON(*data, http.StatusOK)
}

// GetDiscount godoc
// @Tags Discounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} appmodels.DiscountAdminOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/discount/{id} [get]
func GetDiscount(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx)

	baseDB = baseDB.Table("discount d").
		Joins("LEFT JOIN product_item pi2 ON pi2.id = d.product_item_id").
		Joins(`INNER JOIN "user" cu ON cu.id = d.created_by_id`).
		Joins(`LEFT JOIN "product" p ON p.id = pi2.product_id`).
		Joins(`LEFT JOIN "user" ru ON ru.id = d.related_user_id`)

	var data appmodels.DiscountAdminOutPutModel

	if err := baseDB.Select("d.id, d.created_at, d.updated_at, d.product_item_id, p.name as product_name, d.type, d.value, d.quantity, d.expires_in, d.code, d.created_by_id as creator_user_id, cu.username,d.related_user_id, ru.username as related_user_username").First(&data, "d.id", id).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	return ctx.JSON(data, http.StatusOK)
}

// Create Discount godoc
// @Tags Discounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param Discount   body  appmodels.DiscountReqModel  true  "Discount model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/discount  [post]
func CreateDiscount(ctx *app.HttpContext) error {
	var inputModel appmodels.DiscountReqModel

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

// Edit Discount godoc
// @Tags Discounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Param Discount   body  appmodels.DiscountReqModel  true  "Discount model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/discount/edit/{id}  [post]
func EditDiscount(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	var inputModel appmodels.DiscountReqModel

	err = ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	baseDB := db.MustGormDBConn(ctx)

	var dbModel models.Discount

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

// Delete Discount godoc
// @Tags Discounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/discount/delete/{id}  [post]
func DeleteDiscount(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))
	if err != nil {
		return err
	}

	baseDB := db.MustGormDBConn(ctx).Model(&models.Discount{})

	var dbModel models.Discount

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if baseDB.Delete(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, nil)
	}

	return ctx.QuickResponse(consts.Deleted, http.StatusOK)
}

// Validate Discount godoc
// @Tags Discounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param code  path  string  true  "Discount Code"
// @Success 200 {object} appmodels.ValidateDiscountOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/discount/validate/code/{code} [get]
func ValidateDiscountByCode(ctx *app.HttpContext) error {
	code := ctx.GetPathParam("code")

	user, err := ctx.GetUser()
	if err != nil {
		return errors.NewUnauthorizedError(consts.UnauthorizedError, err)
	}

	baseDB := db.MustGormDBConn(ctx)

	var dbModel models.Discount

	if baseDB.First(&dbModel, "code=?", code).Error != nil {
		return errors.NewRecordNotFoundError(consts.ModelDiscountNotFound, nil)
	}

	if err := dbModel.IsValidToUse(user.ID, nil); err != nil {
		return ctx.QuickDBResponse(err.Error(), appmodels.ValidateDiscountOutPutModel{
			Success: false,
		}, http.StatusOK)
	}

	return ctx.QuickDBResponse(consts.OperationDone, appmodels.ValidateDiscountOutPutModel{
		Success: true,
		ID:      dbModel.ID,
		Value:   &dbModel.Value,
		Type:    &dbModel.Type,
	}, http.StatusOK)
}

// Get User Discounts godoc
// @Tags Discounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param searchTerm  query  string  false  "search for item"
// @Success 200 {object} parameter.ListResponse[appmodels.DiscountUserOutPutModel]
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/discount [get]
func GetUserDiscounts(ctx *app.HttpContext) error {
	user, err := ctx.GetUser()
	if err != nil {
		return errors.NewUnauthorizedError(consts.UnauthorizedError, err)
	}
	baseDB := db.MustGormDBConn(ctx).Debug()
	parameter := parameter.New[appmodels.DiscountUserOutPutModel](ctx, baseDB)

	baseDB = baseDB.Table("discount d").
		Joins(`INNER JOIN "user" ru ON ru.id = d.related_user_id`).
		Where("d.deleted_at IS NULL AND ru.id = ?", *user.ID).
		Where("(d.expires_in IS NOT NULL AND d.expires_in >= NOW()) AND (d.quantity IS NOT NULL AND d.quantity > 0)")

	if productItemID, ok := ctx.GetParam("productItemId"); ok {
		baseDB = baseDB.Where("pi2.id = ?", productItemID)
	}

	if relatedUserID, ok := ctx.GetParam("relatedUser"); ok {
		baseDB = baseDB.Where("ru.id = ?", relatedUserID)
	}

	if creatorUserID, ok := ctx.GetParam("createdBy"); ok {
		baseDB = baseDB.Where("cu.id = ?", creatorUserID)
	}

	data, err := parameter.SelectColumns("d.id, d.created_at, d.updated_at, d.type, d.value, d.quantity, d.expires_in, d.code").
		SortDescending("d.created_at", "d.expires_in").
		SearchColumns("d.code").
		EachItemProcess(func(db *gorm.DB, item *appmodels.DiscountUserOutPutModel) error {
			item.TypeName = item.Type.String()
			return nil
		}).
		Execute(baseDB)

	if err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.JSON(*data, http.StatusOK)
}
