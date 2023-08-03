package controllers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/helpers"
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
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param categoryId  query  string  false  "Category ID"
// @Param brandId  query  string  false  "Brand ID"
// @Param minPrice  query  float64  false  "Min Price"
// @Param maxPrice  query  float64  false  "Max Price"
// @Param searchTerm  query  string  false  "search for product name"
// @Success 200 {object} helpers.ListResponse[appmodels.ProductWithItemOutPutModel]
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /product [get]
func GetProducts(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx)

	page, ok := ctx.GetParam("page")
	if !ok {
		page = "1"
	}
	pageInt, err := strconv.Atoi(page)

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	limit, ok := ctx.GetParam("limit")
	if !ok {
		limit = "25"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	baseDB = baseDB.Table("product as p").
		Joins("CROSS JOIN LATERAL (?) as pi2", baseDB.Table("product_item pi2").
			Select("id, price").
			Where("pi2.quantity > 0 AND pi2.product_id = p.id").
			Order("CASE WHEN p.default_product_item_id IS NULL THEN pi2.bought_quantity WHEN pi2.id = p.default_product_item_id THEN 0 ELSE 1 END").
			Limit(1),
		).
		Joins("INNER JOIN brand b ON b.id = p.brand_id").
		Joins("INNER JOIN category c ON c.id = p.category_id").
		Joins("CROSS JOIN LATERAL (?) as pf", baseDB.Table("product_file_map pf").
			Select("file_id").
			Where("pf.product_id = p.id").Limit(1),
		).
		Joins("INNER JOIN file f ON f.id = pf.file_id").
		Where("p.deleted_at IS NULL")

	if categoryID, ok := ctx.GetParam("categoryId"); ok {
		baseDB = baseDB.Where("c.id = ?", categoryID)
	}

	if brandID, ok := ctx.GetParam("brandId"); ok {
		baseDB = baseDB.Where("b.id = ?", brandID)
	}

	if minPrice, ok := ctx.GetParam("minPrice"); ok {
		baseDB = baseDB.Where("pi2.price >= ? ", minPrice)
	}

	if maxPrice, ok := ctx.GetParam("maxPrice"); ok {
		baseDB = baseDB.Where("pi2.price <= ?", maxPrice)
	}

	if searchTerm, ok := ctx.GetParam("searchTerm"); ok {
		baseDB = baseDB.Where("p.name LIKE ?", "%"+strings.TrimSpace(searchTerm)+"%")
	}

	var total int64
	var data []appmodels.ProductWithItemOutPutModel

	var wg sync.WaitGroup
	wg.Add(2)
	errChan := make(chan error, 2)

	go func() {
		defer wg.Done()

		db := baseDB.WithContext(context.Background())
		if err := db.Count(&total).Error; err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()

		db := baseDB.WithContext(context.Background())
		db = db.Offset(limitInt * (pageInt - 1)).Limit(limitInt)
		if err := db.Select("p.id, p.name, p.code, pi2.price, p.brand_id, b.name as brand_name, p.category_id, c.name as category_name, pi2.id as item_id, f.file_type, f.unique_file_name as file_name").Find(&data).Error; err != nil {
			errChan <- err
		}
	}()

	wg.Wait()

	select {
	case err := <-errChan:
		{
			return errors.NewBadRequestError(consts.BadRequest, err)
		}
	default:
		{

		}
	}

	for i := 0; i < len(data); i++ {
		product := data[i]
		data[i].FileUrl = product.FileType.GetDirectory() + "/" + product.FileName
	}

	response := helpers.NewListResponse[appmodels.ProductWithItemOutPutModel](int(pageInt), int(limitInt), total, data)
	return ctx.JSON(*response, http.StatusOK)
}

// GetProducts godoc
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} []appmodels.ProductOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /product/list [get]
func GetProductsList(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx)

	var data []appmodels.ProductOutPutModel

	if err := baseDB.Debug().Table("product as p").
		Joins(`INNER JOIN brand b ON b.id = p.brand_id`).
		Joins("INNER JOIN file f on f.id = b.file_id").
		Joins(`INNER JOIN category c ON c.id = p.category_id`).
		Select(`p.id, p."name", p.code, p.brand_id, 
		b."name" AS brand_name, p.category_id, c."name" AS category_name, f.file_type AS brand_file_type, f.unique_file_name AS brand_file_name`).
		Where("p.deleted_at IS NULL").
		Find(&data).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	for i, product := range data {
		product.BrandFileUrl = models.GetFileUrl(product.BrandFileType, product.BrandFileName)
		data[i] = product
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

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx)

	var data appmodels.ProductOutPutModel

	if err := baseDB.Table("product as p").
		Joins(`INNER JOIN brand b ON b.id = p.brand_id`).
		Joins(`INNER JOIN category c ON c.id = p.category_id`).
		Where("p.id = ? AND p.deleted_at IS NULL", id).
		Select(`p.id, p."name", p.code, p.brand_id, 
		b."name" AS brand_name, p.category_id, c."name" AS category_name`).
		First(&data).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	return ctx.JSON(data, http.StatusOK)
}

// Create Product godoc
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Param Product   body  appmodels.ProductReqModel  true  "Product model"
// @Success 200 {object} helpers.SuccessResponse
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
// @Success 200 {object} helpers.SuccessResponse
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
// @Success 200 {object} helpers.SuccessResponse
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
