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

// GetProducts godoc
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param searchTerm  query  string  false  "search for item"
// @Param categoryId  query  string  false  "Category ID"
// @Param brandId  query  string  false  "Brand ID"
// @Param minPrice  query  float64  false  "Min Price"
// @Param maxPrice  query  float64  false  "Max Price"
// @Param order query string false "order by" Enums(newest,topSell,cheap,expersive)
// @Success 200 {object} parameter.ListResponse[appmodels.ProductWithItemOutPutModel]
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/product [get]
func GetUserProducts(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx).Debug()

	parameter := parameter.New[appmodels.ProductWithItemOutPutModel](ctx, baseDB)

	productItemQry := baseDB.Table("product_item pi2").
		Joins("LEFT JOIN (?) as d ON d.product_item_id = pi2.id", baseDB.Table("discount d").
			Where("d.product_item_id IS NOT NULL AND d.deleted_at IS NULL").
			Where("CASE WHEN d.expires_in IS NOT NULL THEN d.expires_in > NOW() WHEN d.quantity IS NOT NULL THEN d.quantity > 0 ELSE TRUE END").
			Where("d.related_user_id IS NULL").
			Select("d.type, d.value, d.product_item_id, d.quantity").
			Limit(1),
		).
		Select("pi2.id, pi2.price, pi2.created_at, pi2.bought_quantity, d.type, d.value, d.quantity as discount_quantity, pi2.quantity").
		Where("pi2.quantity > 0 AND pi2.product_id = p.id AND pi2.deleted_at IS NULL")

	order, _ := ctx.GetParam("order")

	switch order {
	case "newest":
		productItemQry = productItemQry.Order("pi2.created_at DESC")
	case "topSell":
		productItemQry = productItemQry.Order("pi2.bought_quantity DESC")
	case "cheap":
		productItemQry = productItemQry.Order("pi2.price ASC")
	case "expersive":
		productItemQry = productItemQry.Order("pi2.price DESC")
	default:
		productItemQry = productItemQry.Order("CASE WHEN p.default_product_item_id IS NULL THEN pi2.bought_quantity WHEN pi2.id = p.default_product_item_id THEN 0 ELSE 1 END")
	}

	baseDB = baseDB.Table("product as p").
		Joins("CROSS JOIN LATERAL (?) as pi2", productItemQry.Limit(1)).
		Joins("INNER JOIN brand b ON b.id = p.brand_id").
		Joins("INNER JOIN category c ON c.id = p.category_id").
		Joins("CROSS JOIN LATERAL (?) as pf", baseDB.Table("product_file_map pf").
			Select("file_id").
			Where("pf.product_id = p.id").Order("pf.priority ASC").Limit(1),
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

	switch order {
	case "newest":
		baseDB = baseDB.Order("pi2.created_at DESC")
	case "cheap":
		baseDB = baseDB.Order("pi2.price ASC")
	case "expersive":
		baseDB = baseDB.Order("pi2.price DESC")
	default:
		baseDB = baseDB.Order("pi2.bought_quantity DESC")
	}

	response, err := parameter.SelectColumns("p.id, p.name, p.code, pi2.price, p.brand_id, b.name as brand_name, p.category_id, c.name as category_name, pi2.id as item_id, f.file_type, f.unique_file_name as file_name, pi2.type as discount_type, pi2.value as discount_value, pi2.discount_quantity, pi2.quantity").
		SearchColumns("p.name", "p.code").
		EachItemProcess(func(db *gorm.DB, item *appmodels.ProductWithItemOutPutModel) error {
			item.FileUrl = item.FileType.GetDirectory() + "/" + item.FileName

			return nil
		}).
		Execute(baseDB)

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

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
// @Router /admin/product [get]
func GetAdminProducts(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx)

	var data []appmodels.ProductOutPutModel

	if err := baseDB.Table("product as p").
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
		product.BrandFileUrl = product.BrandFileType.GetFileUrl(product.BrandFileName)
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
// @Router /user/product/{id} [get]
func GetProduct(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx)

	var data appmodels.ProductOutPutModel

	if err := baseDB.Debug().Table("product as p").
		Joins(`INNER JOIN brand b ON b.id = p.brand_id`).
		Joins("INNER JOIN file f on f.id = b.file_id").
		Joins(`INNER JOIN category c ON c.id = p.category_id`).
		Select(`p.id, p."name", p.code, p.brand_id, b."name" AS brand_name, 
			p.category_id, c."name" AS category_name, f.file_type AS brand_file_type, f.unique_file_name AS brand_file_name`).
		Where("p.id = ? AND p.deleted_at IS NULL", id).
		First(&data).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	data.BrandFileUrl = data.BrandFileType.GetFileUrl(data.BrandFileName)

	return ctx.JSON(data, http.StatusOK)
}

// Create Product godoc
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Param Product   body  appmodels.ProductReqModel  true  "Product model"
// @Success 200 {object} helpers.SuccessDBResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/product  [post]
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

	dbModel := inputModel.ToDBModel()
	if err := baseDB.Create(dbModel).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.QuickDBResponse(consts.Created, *dbModel.ID, http.StatusOK)
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
// @Router /admin/product/edit/{id}  [post]
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
// @Router /admin/product/delete/{id}  [post]
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

// Get Suggestion Products godoc
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param searchTerm  query  string  false  "search for item"
// @Param categoryId  query  string  false  "Category ID"
// @Param brandId  query  string  false  "Brand ID"
// @Success 200 {object} parameter.ListResponse[appmodels.SuggestionProductOutPutModel]
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/product/suggestions [get]
func GetSuggestionProducts(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx)
	parameter := parameter.New[appmodels.SuggestionProductOutPutModel](ctx, baseDB)

	qry := baseDB.Table("product as p").
		Joins("CROSS JOIN LATERAL (?) as pi2", baseDB.Table("product_item pi2").
			Select("id, price, color_id").
			Where("pi2.quantity > 0 AND pi2.product_id = p.id").
			Order("CASE WHEN p.default_product_item_id IS NULL THEN pi2.bought_quantity WHEN pi2.id = p.default_product_item_id THEN 0 ELSE 1 END").
			Limit(1),
		).
		Where("p.deleted_at IS NULL AND p.top_features IS NOT NULL AND jsonb_array_length(p.top_features) > 0").
		Where("EXISTS (SELECT true FROM product_file_map WHERE product_id = p.id)")

	if categoryID, ok := ctx.GetParam("categoryId"); ok {
		qry = qry.Where("c.id = ?", categoryID)
	}

	if brandID, ok := ctx.GetParam("brandId"); ok {
		qry = qry.Where("b.id = ?", brandID)
	}

	response, err := parameter.SelectColumns("p.id as product_id", "p.name", "pi2.id as product_item_id", "pi2.color_id", "p.top_features").
		SearchColumns("p.name", "p.code").
		EachItemProcess(func(db *gorm.DB, data *appmodels.SuggestionProductOutPutModel) error {

			// files
			var files []appmodels.ProductItemFileOutPutModel

			if err := baseDB.Table("file as f").
				Joins("INNER JOIN product_file_map pf ON pf.file_id = f.id").
				Where("pf.product_id = ?", data.ProductID).Find(&files).Error; err != nil {
				return errors.NewInternalServerError(consts.InternalServerError, nil)
			}

			for i := 0; i < len(files); i++ {
				file := files[i]
				files[i].FileUrl = file.FileType.GetDirectory() + "/" + file.UniqueFileName
			}

			data.Files = files

			// colors
			var colors []appmodels.ProductItemInfoColorOutPutModel

			if err := baseDB.Table("product_item pi2").
				Joins("INNER JOIN color c ON c.id = pi2.color_id").
				Where("pi2.deleted_at IS NULL AND pi2.product_id=?", data.ProductID).
				Select("pi2.id AS product_item_id, c.name, c.color_hex").Find(&colors).Error; err != nil {
				return errors.NewInternalServerError(consts.InternalServerError, err)
			}

			data.Colors = colors

			return nil
		}).Execute(qry)

	if err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.JSON(*response, http.StatusOK)
}

// GetProduct godoc
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} appmodels.ProductAdminOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/product/{id} [get]
func GetAdminProduct(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx)

	var data appmodels.ProductAdminOutPutModel

	if err := baseDB.Table("product as p").
		Joins(`INNER JOIN brand b ON b.id = p.brand_id`).
		Joins("INNER JOIN file f on f.id = b.file_id").
		Joins(`INNER JOIN category c ON c.id = p.category_id`).
		Select(`p.id, p."name", p.code, p.brand_id, b."name" AS brand_name, 
			p.category_id, c."name" AS category_name, p.description, p.short_description, p.top_features`).
		Where("p.id = ? AND p.deleted_at IS NULL", id).
		First(&data, id).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	return ctx.JSON(data, http.StatusOK)
}
