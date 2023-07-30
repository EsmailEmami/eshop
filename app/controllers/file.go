package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/consts"
	dbpkg "github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/errors"
	"github.com/esmailemami/eshop/models"
	service "github.com/esmailemami/eshop/services/file"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UploadFile godoc
// @Summary Upload an image file
// @Description Uploads an image file to the server.
// @Tags Files
// @Accept json
// @Produce json
// @Security Bearer
// @Param itemId  path  string  true  "item ID"
// @Param fileType  path  int  true  "file Type"
// @Param file formData file true "Image file to be uploaded"
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /file/uploadImage/{itemId}/{fileType} [post]
func UploadImage(ctx *app.HttpContext) error {
	err := ctx.Request.ParseMultipartForm(10 << 20) // 10 MB maximum file size
	if err != nil {
		return errors.NewBadRequestError(consts.InvalidFileSize, err)
	}

	fileTypeInput, err := strconv.Atoi(ctx.GetPathParam("fileType"))
	if err != nil {
		return errors.NewBadRequestError("invalid file type", err)
	}

	itemID, err := uuid.Parse(ctx.GetPathParam("itemId"))
	if err != nil {
		return errors.NewBadRequestError("invalid itemId", err)
	}

	fileType, err := models.FileTypeFromInt(fileTypeInput)
	if err != nil {
		return errors.NewBadRequestError(err.Error(), err)
	}

	baseDB := dbpkg.MustGormDBConn(ctx)
	baseTx := baseDB.Begin()

	multiple, err := validateItem(baseDB, itemID, fileType)

	if err != nil {
		return errors.NewBadRequestError(err.Error(), err)
	}

	paths := []string{}
	files := []*models.File{}

	errCh := make(chan struct{})

	go func() {
		<-errCh
		for _, path := range paths {
			service.DeleteFile(path)
		}
	}()

	// Process each part of the multipart request
loop:
	for _, fileHeaders := range ctx.Request.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			if !service.IsImageFile(fileHeader) {
				fmt.Println("this is not image!")
			}
			path, fileName, err := service.UploadFile(fileHeader, "uploads/test", true, true)

			if err != nil {
				errCh <- struct{}{}
				return err
			}
			paths = append(paths, path)

			files = append(files, &models.File{
				Model: models.Model{
					ID: models.NewID(),
				},
				MimeType:       service.GetMimeType(fileHeader),
				Extension:      service.GetFileExetension(fileName),
				OriginalName:   fileHeader.Filename,
				UniqueFileName: fileName,
				FileType:       fileType,
			})

			if !multiple {
				break loop
			}
		}
	}

	if len(files) == 0 {
		return ctx.QuickResponse("فایلی به سرور ارسال نشده است.", http.StatusBadRequest)
	}

	if err := baseTx.CreateInBatches(files, len(files)).Error; err != nil {
		baseTx.Rollback()
		errCh <- struct{}{}
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	err = insertItemFile(baseDB, baseTx, itemID, fileType, files...)

	if err != nil {
		baseTx.Rollback()
		errCh <- struct{}{}
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	baseTx.Commit()

	return ctx.QuickResponse("عملیات با موفقیت به پایان رسید", http.StatusOK)
}

func validateItem(db *gorm.DB, itemID uuid.UUID, fileType models.FileType) (multiple bool, err error) {
	switch fileType {
	case models.FileTypeSystematic:
		return true, nil
	case models.FileTypeProduct:
		if !dbpkg.Exists(db, &models.Product{}, "id = ?", itemID) {
			return true, fmt.Errorf("no item found with Id #%s", itemID.String())
		}

		return true, nil
	case models.FileTypeBrand:
		if !dbpkg.Exists(db, &models.Brand{}, "id = ?", itemID) {
			return false, fmt.Errorf("no item found with Id #%s", itemID.String())
		}

		return false, nil
	default:
		return true, fmt.Errorf("invalid file type")
	}
}

func insertItemFile(db, tx *gorm.DB, itemID uuid.UUID, fileType models.FileType, files ...*models.File) error {
	switch fileType {
	case models.FileTypeSystematic:
		return nil
	case models.FileTypeProduct:
		mapItems := []models.ProductFileMap{}

		for _, file := range files {
			mapItems = append(mapItems, models.ProductFileMap{
				ProductID: itemID,
				FileID:    *file.ID,
			})
		}

		return tx.CreateInBatches(mapItems, len(mapItems)).Error
	case models.FileTypeBrand:

		return tx.Model(&models.Brand{}).Where("id = ?", itemID).UpdateColumn("file_id", *files[0].ID).Error
	default:
		return nil
	}
}
