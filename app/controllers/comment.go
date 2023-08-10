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

// Get Admin Comments godoc
// @Tags Comments
// @Accept json
// @Produce json
// @Security Bearer
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param searchTerm  query  string  false  "search for item"
// @Param productId  query  string  false  "Product ID"
// @Param userId  query  string  false  "User ID"
// @Success 200 {object} parameter.ListResponse[appmodels.CommentOutPutModel]
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/comment [get]
func GetAdminUserComments(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx).Table(`"comment" c`).
		Joins(`INNER JOIN "user" u ON u.id = c.created_by_id`).
		Where("c.deleted_at IS NULL")

	if productID, ok := ctx.GetParam("productId"); ok {
		baseDB = baseDB.Where("c.product_id=?", productID)
	}

	if userID, ok := ctx.GetParam("userId"); ok {
		baseDB = baseDB.Where("c.created_by_id=?", userID)
	}

	parameter := parameter.New[appmodels.CommentOutPutModel](ctx, baseDB)

	response, err := parameter.SearchColumns("c.id, c.created_at, c.updated_at, c.text,c.rate,c.strength_points,c.weak_ponits, u.username").
		SortDescending("c.created_at").Execute(baseDB)

	if err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.JSON(*response, http.StatusOK)
}

// Get User Comments godoc
// @Tags Comments
// @Accept json
// @Produce json
// @Security Bearer
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param searchTerm  query  string  false  "search for item"
// @Success 200 {object} parameter.ListResponse[appmodels.CommentOutPutModel]
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/comment [get]
func GetUserComments(ctx *app.HttpContext) error {
	user, err := ctx.GetUser()

	if err != nil {
		return errors.NewUnauthorizedError(consts.UnauthorizedError, err)
	}

	baseDB := db.MustGormDBConn(ctx).Table(`"comment" c`).
		Joins(`INNER JOIN "user" u ON u.id = c.created_by_id`).
		Where("c.created_by_id=? AND c.deleted_at IS NULL", *user.ID)

	parameter := parameter.New[appmodels.CommentOutPutModel](ctx, baseDB)

	response, err := parameter.SearchColumns("c.id, c.created_at, c.updated_at, c.text,c.rate,c.strength_points,c.weak_ponits, u.username").
		SortDescending("c.created_at").Execute(baseDB)

	if err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.JSON(*response, http.StatusOK)
}

// Get Product Comments godoc
// @Tags Comments
// @Accept json
// @Produce json
// @Security Bearer
// @Param productId  path  string  true  "Product ID"
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param searchTerm  query  string  false  "search for item"
// @Success 200 {object} parameter.ListResponse[appmodels.CommentOutPutModel]
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/comment/product/{productId} [get]
func GetProductComments(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx).Table(`"comment" c`).
		Joins(`INNER JOIN "user" u ON u.id = c.created_by_id`).
		Where("c.deleted_at IS NULL")

	productID, err := uuid.Parse(ctx.GetPathParam("productId"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	baseDB = baseDB.Where("c.product_id=?", productID)

	parameter := parameter.New[appmodels.CommentOutPutModel](ctx, baseDB)

	response, err := parameter.SearchColumns("c.id, c.created_at, c.updated_at, c.text,c.rate,c.strength_points,c.weak_ponits, u.username").
		SortDescending("c.created_at").Execute(baseDB)

	if err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.JSON(*response, http.StatusOK)
}

// GetComment godoc
// @Tags Comments
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} appmodels.CommentOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/comment/{id} [get]
func GetComment(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx)

	var data appmodels.CommentOutPutModel

	if err := baseDB.Table(`"comment" c`).
		Joins(`INNER JOIN "user" u ON u.id = c.created_by_id`).
		Select("c.id, c.created_at, c.updated_at, c.text,c.rate,c.strength_points,c.weak_ponits, u.username").
		First(&data, "c.id", id).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	return ctx.JSON(data, http.StatusOK)
}

// Create Comment godoc
// @Tags Comments
// @Accept json
// @Produce json
// @Security Bearer
// @Param Comment   body  appmodels.CommentReqModel  true  "Comment model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/comment  [post]
func CreateComment(ctx *app.HttpContext) error {
	var inputModel appmodels.CommentReqModel

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

// Edit Comment godoc
// @Tags Comments
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Param Comment   body  appmodels.CommentReqModel  true  "Comment model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/comment/edit/{id}  [post]
func EditComment(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	var inputModel appmodels.CommentReqModel

	err = ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx)

	err = inputModel.ValidateUpdate(baseDB)
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	var dbModel models.Comment

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	inputModel.MergeWithDBData(&dbModel)
	if baseDB.Save(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.QuickResponse(consts.Updated, http.StatusOK)
}

// Delete Comment godoc
// @Tags Comments
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/comment/delete/{id}  [post]
func DeleteComment(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))
	if err != nil {
		return err
	}

	baseDB := db.MustGormDBConn(ctx).Model(&models.Comment{})

	var dbModel models.Comment

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if baseDB.Delete(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, nil)
	}

	return ctx.QuickResponse(consts.Deleted, http.StatusOK)
}
