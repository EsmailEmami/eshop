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

// GetProducts godoc
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} []appmodels.ProductOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /product [get]
func GetProducts(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx).Model(&models.Product{})

	var data []appmodels.ProductOutPutModel

	if err := baseDB.Table("product as p").
		Joins(`INNER JOIN product_item pi2 ON pi2.product_id = p.id
	       INNER JOIN brand b ON b.id = p.brand_id
	       INNER JOIN category c ON c.id = p.category_id
	       INNER JOIN product_file_map pf ON pf.product_id = p.id
	       INNER JOIN file f ON f.id = pf.file_id`).
		Order(`p.id, f.created_at,
		  CASE
			  WHEN p.default_product_item_id IS NULL THEN pi2.bought_quantity
			  WHEN pi2.id = p.default_product_item_id THEN 0
			  ELSE 1 
		  END, pi2.bought_quantity`).
		Select(`DISTINCT ON (p.id) p.id, p."name", p.code, pi2.price, p.brand_id, 
	        b."name" AS brand_name, p.category_id, c."name" AS category_name, 
	        pi2.id AS item_id, f.file_type, f.unique_file_name AS file_name`).
		Find(&data).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	for i := 0; i < len(data); i++ {
		product := data[i]
		data[i].FileUrl = product.FileType.GetDirectory() + "/" + product.FileName
	}

	return ctx.JSON(data, http.StatusOK)
}

// GetProduct godoc
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} appmodels.ProductOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /product/{id} [get]
func GetProduct(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	_ = id
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx).Model(&models.Product{})

	var data appmodels.ProductOutPutModel

	if err := baseDB.Table("product as p").
		Joins(`INNER JOIN product_item pi2 ON pi2.product_id = p.id
	   INNER JOIN brand b ON b.id = p.brand_id
	   INNER JOIN category c ON c.id = p.category_id
	   INNER JOIN product_file_map pf ON pf.product_id = p.id
	   INNER JOIN file f ON f.id = pf.file_id`).
		Order(`p.id, f.created_at,
	  CASE
		  WHEN p.default_product_item_id IS NULL THEN pi2.bought_quantity
		  WHEN pi2.id = p.default_product_item_id THEN 0
		  ELSE 1 
	  END, pi2.bought_quantity`).
		Where("p.id = ?", id).
		Select(`DISTINCT ON (p.id) p.id, p."name", p.code, pi2.price, p.brand_id, 
		b."name" AS brand_name, p.category_id, c."name" AS category_name, 
		pi2.id AS item_id, f.file_type, f.unique_file_name AS file_name`).
		First(&data).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	data.FileUrl = data.FileType.GetDirectory() + "/" + data.FileName
	return ctx.JSON(data, http.StatusOK)
}

// Create Product godoc
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Param Product   body  appmodels.ProductReqModel  true  "Product model"
// @Success 200 {object} httpmodels.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /product  [post]
func CreateProduct(ctx *app.HttpContext) error {
	var inputModel appmodels.ProductReqModel

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

// Edit Product godoc
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Param Product   body  appmodels.ProductReqModel  true  "Product model"
// @Success 200 {object} httpmodels.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /product/edit/{id}  [post]
func EditProduct(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	var inputModel appmodels.ProductReqModel

	err = ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	baseDB := db.MustGormDBConn(ctx)

	err = inputModel.ValidateUpdate(baseDB)
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	var dbModel models.Product

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if db.Exists(baseDB, &models.Product{}, "code = ? and id != ?", inputModel.Code, id) {
		return errors.NewValidationError(consts.ExistedCode, nil)
	}

	inputModel.MergeWithDBData(&dbModel)
	if baseDB.Save(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.QuickResponse(consts.Updated, http.StatusOK)
}

// Delete Product godoc
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} httpmodels.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /product/delete/{id}  [post]
func DeleteProduct(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))
	if err != nil {
		return err
	}

	baseDB := db.MustGormDBConn(ctx).Model(&models.Product{})

	var dbModel models.Product

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if baseDB.Delete(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, nil)
	}

	return ctx.QuickResponse(consts.Deleted, http.StatusOK)
}
