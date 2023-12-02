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
)

// GetProductItem godoc
// @Tags ProductItems
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} appmodels.ProductItemInfoOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/productItem/{id} [get]
func GetProductItem(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx)

	var data appmodels.ProductItemInfoOutPutModel

	if err := baseDB.Table("product_item pi2").
		Joins("INNER JOIN product p ON P.id = pi2 .product_id").
		Joins("INNER JOIN color c ON C.id = pi2.color_id").
		Joins("LEFT JOIN LATERAL (?) as d ON TRUE", baseDB.Table("discount d").
			Where("d.product_item_id = pi2.id AND d.deleted_at IS NULL AND d.type = 1").
			Where("(d.expires_in IS NOT NULL AND d.expires_in >= NOW()) AND (d.quantity IS NOT NULL AND d.quantity > 0)").
			Where("d.related_user_id IS NULL").
			Order("d.created_at ASC").
			Select("d.type, d.value, d.product_item_id, d.quantity").
			Limit(1),
		).
		Select(`pi2.id, pi2.price,pi2.status, pi2 .color_id, pi2.product_id, pi2.quantity,
		p."name" AS product_title,p.rate, p.code AS product_code, p.short_description AS product_short_description, 
		p.description AS product_description, c."name" AS color_name,d.type as discount_type, d.value as discount_value, d.quantity as discount_quantity`).
		First(&data, "pi2.id", id).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	// files query
	var files []appmodels.ProductItemFileOutPutModel

	if err := baseDB.Table("file as f").
		Joins("INNER JOIN product_file_map pf ON pf.file_id = f.id").
		Where("pf.product_id = ?", data.ProductID).Find(&files).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, nil)
	}

	for i := 0; i < len(files); i++ {
		file := files[i]
		files[i].FileUrl = file.FileType.GetFileUrl(file.UniqueFileName)
	}

	data.Files = files

	// features
	data.Features = []appmodels.ProductItemCategoryFeatureModel{}

	featureData := []struct {
		Category string `gorm:"column:category"`
		Key      string `gorm:"column:key"`
		Value    string `gorm:"column:value"`
	}{}

	if err := baseDB.Table("product_feature_value pfv").
		Joins("INNER JOIN product_feature_key pfk ON pfk.id = pfv.product_feature_key_id").
		Joins("INNER JOIN product_feature_category pfc ON pfc.id = pfk.product_feature_category_id").
		Where("pfv.product_id = ?", data.ProductID).
		Select("pfk.name as key,pfv.value, pfc.name as category").Find(&featureData).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	features := map[string][]appmodels.ProductItemFeatureModel{}

	for _, v := range featureData {
		data, exists := features[v.Category]

		if exists {
			data = append(data, appmodels.ProductItemFeatureModel{
				Key:   v.Key,
				Value: v.Value,
			})
		} else {
			data = []appmodels.ProductItemFeatureModel{
				{
					Key:   v.Key,
					Value: v.Value,
				},
			}
		}

		features[v.Category] = data
	}

	for categoryName, keys := range features {
		data.Features = append(data.Features, appmodels.ProductItemCategoryFeatureModel{
			Category: categoryName,
			Items:    keys,
		})
	}

	// colors
	var colors []appmodels.ProductItemInfoColorOutPutModel

	if err := baseDB.Table("product_item pi2").
		Joins("INNER JOIN color c ON c.id = pi2.color_id").
		Where("pi2.deleted_at IS NULL AND pi2.product_id=?", data.ProductID).
		Select("pi2.id AS product_item_id, c.name, c.color_hex").Find(&colors).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	data.Colors = colors

	return ctx.JSON(data, http.StatusOK)
}

// GetProductItem godoc
// @Tags ProductItems
// @Accept json
// @Produce json
// @Security Bearer
// @Param productId  path  string  true  "Product ID"
// @Success 200 {object} []appmodels.ProductOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/productItem/product/{productId} [get]
func GetProductItems(ctx *app.HttpContext) error {
	productID, err := uuid.Parse(ctx.GetPathParam("productId"))
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx)

	var data []appmodels.ProductItemOutPutModel

	if err := baseDB.Table("product_item pi2").
		Joins("INNER JOIN product p ON P.id = pi2 .product_id").
		Joins("INNER JOIN color c ON C.id = pi2.color_id").
		Select(`pi2.id, pi2.price,pi2.status, pi2 .color_id, pi2.product_id, pi2.quantity,
		p."name" AS product_title, p.code AS product_code, c."name" AS color_name`).
		Find(&data, "p.id = ? AND pi2.deleted_at IS NULL", productID).Error; err != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	return ctx.JSON(data, http.StatusOK)
}

// Create ProductItem godoc
// @Tags ProductItems
// @Accept json
// @Produce json
// @Security Bearer
// @Param ProductItem   body  appmodels.ProductItemReqModel  true  "ProductItem model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/productItem  [post]
func CreateProductItem(ctx *app.HttpContext) error {
	var inputModel appmodels.ProductItemReqModel

	err := ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx)
	baseTx := baseDB.Begin()

	err = inputModel.ValidateCreate()
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	dbModel := inputModel.ToDBModel()

	if db.Exists(
		baseDB,
		&models.ProductItem{},
		"color_id=? AND product_id=?",
		dbModel.ColorID,
		dbModel.ProductID,
	) {
		return errors.NewBadRequestError(
			"The item with the entered color has been previously registered.",
			nil,
		)
	}

	if err := baseTx.Create(dbModel).Error; err != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	if inputModel.IsMainItem {
		if err := baseTx.Model(&models.Product{}).
			Where("id=?", inputModel.ProductID).
			UpdateColumn("default_product_item_id", dbModel.ID).Error; err != nil {
			baseTx.Rollback()
			return errors.NewInternalServerError(consts.InternalServerError, err)
		}
	}

	baseTx.Commit()

	return ctx.QuickResponse(consts.Created, http.StatusOK)
}

// Edit ProductItem godoc
// @Tags ProductItems
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Param ProductItem   body  appmodels.ProductItemReqModel  true  "ProductItem model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/productItem/edit/{id}  [post]
func EditProductItem(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	var inputModel appmodels.ProductItemReqModel

	err = ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	baseDB := db.MustGormDBConn(ctx)
	baseTx := baseDB.Begin()

	var dbModel models.ProductItem

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	err = inputModel.ValidateUpdate()
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	inputModel.MergeWithDBData(&dbModel)

	if db.Exists(
		baseDB,
		&models.ProductItem{},
		"color_id=? AND product_id=? AND id!=?",
		dbModel.ColorID,
		dbModel.ProductID,
		*dbModel.ID,
	) {
		return errors.NewBadRequestError(
			"The item with the entered color has been previously registered.",
			nil,
		)
	}

	if baseTx.Save(&dbModel).Error != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	if inputModel.IsMainItem {
		if err := baseTx.Model(&models.Product{}).
			Where("id=?", inputModel.ProductID).
			UpdateColumn("default_product_item_id", dbModel.ID).Error; err != nil {
			baseTx.Rollback()
			return errors.NewInternalServerError(consts.InternalServerError, err)
		}
	}

	baseTx.Commit()

	return ctx.QuickResponse(consts.Updated, http.StatusOK)
}

// Delete ProductItem godoc
// @Tags ProductItems
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/productItem/delete/{id}  [post]
func DeleteProductItem(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))
	if err != nil {
		return err
	}

	baseDB := db.MustGormDBConn(ctx)
	baseTx := baseDB.Begin()

	var dbModel models.ProductItem

	if baseDB.Preload("Product").First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if baseTx.Where("id=?", &dbModel.ID).Delete(&dbModel).Error != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, nil)
	}

	// choose last most sell item
	if dbModel.Product.DefaultProductItemID.String() == dbModel.ID.String() {
		mostSellDBItem := baseDB.Model(&models.ProductItem{}).
			Where("id != ?", dbModel.ID).Order("bought_quantity DESC").
			Select("TOP(1) id")

		if err := baseTx.Model(&models.Product{}).
			Where("id=?", dbModel.ProductID).
			UpdateColumn("default_product_item_id", mostSellDBItem).Error; err != nil {
			baseTx.Rollback()
			return errors.NewInternalServerError(consts.InternalServerError, err)
		}
	}

	baseTx.Commit()

	return ctx.QuickResponse(consts.Deleted, http.StatusOK)
}

// GetProductItemsSelectList godoc
// @Tags ProductItems
// @Accept json
// @Produce json
// @Security Bearer
// @Param productId  path  string  true  "Product ID"
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param searchTerm  query  string  false  "search for item"
// @Success 200 {object} parameter.ListResponse[appmodels.ProductItemAdminSelectListOutPutModel]
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/productItem/selectList/{productId} [get]
func GetProductItemsSelectList(ctx *app.HttpContext) error {
	productID, err := uuid.Parse(ctx.GetPathParam("productId"))
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	baseDB := db.MustGormDBConn(ctx)

	parameter := parameter.New[appmodels.ProductItemAdminSelectListOutPutModel](ctx, baseDB)

	baseDB = baseDB.Table("product_item as pi2").
		Joins("INNER JOIN color c ON c.id = pi2.color_id").
		Where("pi2.deleted_at IS NULL AND pi2.product_id = ?", productID)

	response, err := parameter.SelectColumns("pi2.id, pi2.price, c.id as color_id, c.name as color_name").
		SearchColumns("c.name", "c.code", "c.color_hex").
		Execute(baseDB)

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	return ctx.JSON(*response, http.StatusOK)
}
