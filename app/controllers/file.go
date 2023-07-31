package controllers

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/esmailemami/eshop/app"
	appmodels "github.com/esmailemami/eshop/app/models"
	"github.com/esmailemami/eshop/consts"
	dbpkg "github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/errors"
	"github.com/esmailemami/eshop/models"
	service "github.com/esmailemami/eshop/services/file"
	"github.com/google/uuid"
)

// UploadFile godoc
// @Summary Upload an image file
// @Description Uploads an image file to the server.
// @Tags Files
// @Accept json
// @Produce json
// @Security Bearer
// @Param itemId  path  string  true  "item ID"
// @Param fileType  path  models.FileType  true  "file Type"
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

	multiple, err := service.ValidateItem(baseDB, itemID, fileType)

	if err != nil {
		return errors.NewBadRequestError(err.Error(), err)
	}

	paths := []string{}
	files := []*models.File{}

	errCh := make(chan struct{})

	go func() {
		<-errCh
		for _, path := range paths {
			service.DeleteFileByPath(path)
		}
	}()

loop:
	for _, fileHeaders := range ctx.Request.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			if !service.IsImageFile(fileHeader) {
				fmt.Println("this is not image!")
			}
			path, fileName, err := service.UploadFile(fileHeader, fileType.GetDirectory(), true, true)

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

	err = service.InsertItemFile(baseDB, baseTx, itemID, fileType, files...)

	if err != nil {
		baseTx.Rollback()
		errCh <- struct{}{}
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	baseTx.Commit()

	return ctx.QuickResponse("عملیات با موفقیت به پایان رسید", http.StatusOK)
}

// DeleteFile godoc
// @Tags Files
// @Accept json
// @Produce json
// @Security Bearer
// @Param fileId  path  string  true  "file ID"
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /file/delete/{fileId} [post]
func DeleteFile(ctx *app.HttpContext) error {
	fileID, err := uuid.Parse(ctx.GetPathParam("fileId"))
	if err != nil {
		return errors.NewBadRequestError("invalid fileId", err)
	}
	var file models.File

	baseDB := dbpkg.MustGormDBConn(ctx)

	if baseDB.First(&file, fileID).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	if err := service.DeleteFile(baseDB, &file); err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.QuickResponse(consts.Deleted, http.StatusOK)
}

// GetFile godoc
// @Tags Files
// @Accept json
// @Produce json
// @Security Bearer
// @Param fileId  path  string  true  "file ID"
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /file/{fileId} [get]
func GetFile(ctx *app.HttpContext) error {
	fileID, err := uuid.Parse(ctx.GetPathParam("fileId"))
	if err != nil {
		return errors.NewBadRequestError("invalid fileId", err)
	}
	var dbFile models.File

	baseDB := dbpkg.MustGormDBConn(ctx)

	if baseDB.First(&dbFile, fileID).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	path := service.GetFilePath(&dbFile)

	fileInfo, err := os.Stat(path)
	if err != nil {
		return errors.NewRecordNotFoundError(consts.FileNotFound, err)
	}

	file, err := os.Open(path)
	if err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}
	defer file.Close()

	bts, err := io.ReadAll(file)
	if err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	ctx.ResponseWriter.Header().Set("Content-Type", dbFile.MimeType)
	ctx.ResponseWriter.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	ctx.ResponseWriter.Write(bts)

	return nil
}

// GetStreamingFile godoc
// @Tags Files
// @Accept json
// @Produce json
// @Security Bearer
// @Param fileId  path  string  true  "file ID"
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /file/stream/{fileId} [get]
func GetStreamingFile(ctx *app.HttpContext) error {
	fileID, err := uuid.Parse(ctx.GetPathParam("fileId"))
	if err != nil {
		return errors.NewBadRequestError("invalid fileId", err)
	}
	var dbFile models.File

	baseDB := dbpkg.MustGormDBConn(ctx)

	if baseDB.First(&dbFile, fileID).Error != nil {
		return errors.NewRecordNotFoundError(consts.RecordNotFound, nil)
	}

	path := service.GetFilePath(&dbFile)

	fileInfo, err := os.Stat(path)
	if err != nil {
		return errors.NewRecordNotFoundError(consts.FileNotFound, err)
	}

	file, err := os.Open(path)
	if err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}
	defer file.Close()

	ctx.ResponseWriter.Header().Set("Content-Type", dbFile.MimeType)
	ctx.ResponseWriter.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	bufferedWriter := bufio.NewWriter(ctx.ResponseWriter)
	defer bufferedWriter.Flush()

	const chunkSize = 10 * 1024
	buffer := make([]byte, chunkSize)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			break
		}
		_, err = bufferedWriter.Write(buffer[:n])
		if err != nil {
			break
		}
	}

	if err != nil && err != io.EOF {
		return errors.NewInternalServerError("Failed to stream file content to response", err)
	}

	return nil
}

// GetItemFiles godoc
// @Tags Files
// @Accept json
// @Produce json
// @Security Bearer
// @Param itemId  path  string  true  "item ID"
// @Param fileType  path  models.FileType  true  "file Type"
// @Success 200 {object} []appmodels.FileOutPutModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /file/{itemId}/{fileType} [get]
func GetItemFiles(ctx *app.HttpContext) error {
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

	_, err = service.ValidateItem(baseDB, itemID, fileType)

	if err != nil {
		return errors.NewBadRequestError(err.Error(), err)
	}

	whereClause, ok := service.GenrateFileWhereClause(baseDB, itemID, fileType)

	if !ok {
		return errors.NewBadRequestError("invalid item for get files", nil)
	}

	var files []appmodels.FileOutPutModel

	err = baseDB.Model(&models.File{}).Where("id in (?)", whereClause).Find(&files).Error
	if err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.JSON(files, http.StatusOK)
}
