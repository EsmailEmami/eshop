package controllers

import (
	"fmt"
	"strconv"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/consts"
	"github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/errors"
	"github.com/esmailemami/eshop/models"
	service "github.com/esmailemami/eshop/services/file"
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

	fileType, err := models.FileTypeFromInt(fileTypeInput)
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
		}
	}

	baseDB := db.MustGormDBConn(ctx)
	baseTx := baseDB.Begin()

	if err := baseTx.CreateInBatches(files, len(files)).Error; err != nil {
		baseTx.Rollback()
		errCh <- struct{}{}
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return nil
}
