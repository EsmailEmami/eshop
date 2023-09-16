package controllers

import (
	"net/http"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/errors"
	"github.com/esmailemami/eshop/app/helpers"
	appmodels "github.com/esmailemami/eshop/app/models"
	"github.com/esmailemami/eshop/app/parameter"
	fileService "github.com/esmailemami/eshop/app/services/file"
	"github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetBrands godoc
// @Tags Brands
// @Accept json
// @Produce json
// @Security Bearer
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param searchTerm  query  string  false  "search for item"
// @Success 200 {object} parameter.ListResponse[appmodels.BrandOutPutModel]
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/brand [get]
// @Router /user/brand [get]
func GetBrands(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx)

	parameter := parameter.New[appmodels.BrandOutPutModel](ctx, baseDB)

	baseDB = baseDB.Table("brand b").
		Joins("INNER JOIN file f on f.id = b.file_id").
		Where("b.deleted_at IS NULL")

	response, err := parameter.SelectColumns("b.id, b.created_at, b.updated_at,b.name,b.code, b.file_id, f.unique_file_name as file_name,f.file_type").
		SearchColumns("b.name").
		EachItemProcess(func(db *gorm.DB, t *appmodels.BrandOutPutModel) error {
			t.FileUrl = t.FileType.GetFileUrl(t.FileName)
			return nil
		}).
		SortDescending("b.created_at", "b.name").
		Execute(baseDB)

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	return ctx.JSON(*response, http.StatusOK)
}

// GetBrandsSelectList godoc
// @Tags Brands
// @Accept json
// @Produce json
// @Security Bearer
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param searchTerm  query  string  false  "search for item"
// @Success 200 {object} parameter.ListResponse[helpers.KeyValueResponse[uuid.UUID, string]]
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/brand/selectList [get]
func GetBrandsSelectList(ctx *app.HttpContext) error {
	baseDB := db.MustGormDBConn(ctx)

	parameter := parameter.New[helpers.KeyValueResponse[uuid.UUID, string]](ctx, baseDB)

	baseDB = baseDB.Model(&models.Brand{})

	response, err := parameter.SelectColumns(`id as "key", "name" as "value"`).
		SearchColumns(`name`).
		SortDescending("created_at", `name`).
		Execute(baseDB)

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	return ctx.JSON(*response, http.StatusOK)
}

// GetBrand godoc
// @Tags Brands
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} appmodels.BrandOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/brand/{id} [get]
func GetBrand(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx)

	var brand appmodels.BrandOutPutModel

	if err := baseDB.Table("brand b").
		Joins("INNER JOIN file f on f.id = b.file_id").
		Select("b.id, b.created_at, b.updated_at,b.name,b.code, b.file_id, f.unique_file_name as file_name,f.file_type").
		First(&brand, "b.id", id).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	brand.FileUrl = brand.FileType.GetFileUrl(brand.FileName)

	return ctx.JSON(brand, http.StatusOK)
}

// Create Brand godoc
// @Tags Brands
// @Accept json
// @Produce json
// @Security Bearer
// @Param Brand   body  appmodels.BrandReqModel  true  "Brand model"
// @Success 200 {object} helpers.SuccessDBResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/brand  [post]
func CreateBrand(ctx *app.HttpContext) error {
	var inputModel appmodels.BrandReqModel

	err := ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	baseDB := db.MustGormDBConn(ctx)

	err = inputModel.ValidateCreate()
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	dbModels := inputModel.ToDBModel()
	if err := baseDB.Create(dbModels).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.QuickDBResponse(consts.Created, *dbModels.ID, http.StatusOK)
}

// Edit Brand godoc
// @Tags Brands
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Param Brand   body  appmodels.BrandReqModel  true  "Brand model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/brand/edit/{id}  [post]
func EditBrand(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))

	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	var inputModel appmodels.BrandReqModel

	err = ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	baseDB := db.MustGormDBConn(ctx)
	baseTx := baseDB.Begin()

	var dbModel models.Brand

	if baseDB.First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	err = inputModel.ValidateUpdate(id)
	if err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	inputModel.MergeWithDBData(&dbModel)
	if baseTx.Save(&dbModel).Error != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	baseTx.Commit()

	return ctx.QuickResponse(consts.Updated, http.StatusOK)
}

// Delete Brand godoc
// @Tags Brands
// @Accept json
// @Produce json
// @Security Bearer
// @Param id  path  string  true  "Record ID"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/brand/delete/{id}  [post]
func DeleteBrand(ctx *app.HttpContext) error {
	id, err := uuid.Parse(ctx.GetPathParam("id"))
	if err != nil {
		return err
	}

	baseDB := db.MustGormDBConn(ctx)
	baseTx := baseDB.Begin()

	var dbModel models.Brand

	if baseDB.Preload("File").First(&dbModel, id).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	err = fileService.DeleteFile(baseDB, baseTx, dbModel.ID, dbModel.File)

	if err != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	if baseTx.Delete(&models.Brand{}, dbModel.ID).Error != nil {
		baseTx.Rollback()
		return errors.NewInternalServerError(consts.InternalServerError, nil)
	}

	baseTx.Commit()

	return ctx.QuickResponse(consts.Deleted, http.StatusOK)
}
