package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	Model

	MimeType       string           `gorm:"column:mime_type"                 json:"mimeType"`
	Extension      string           `gorm:"column:extension"                 json:"extension"`
	OriginalName   string           `gorm:"column:original_name"             json:"originalName"`
	UniqueFileName string           `gorm:"column:unique_file_name"          json:"uniqueFineName"`
	FileType       FileType         `gorm:"column:file_type"                 json:"fileType"`
	ItemID         *uuid.UUID       `gorm:"column:item_id"                   json:"itemId"`
	Products       []ProductFileMap `gorm:"foreignKey:file_id;references:id" json:"products"`
	Brands         []Brand          `gorm:"foreignKey:file_id;references:id" json:"brands"`
	AppPics        []AppPic         `gorm:"foreignKey:file_id;references:id" json:"appPics"`
	ExpireDate     *time.Time       `gorm:"expire_date"                      json:"expireDate"`
}

func (File) TableName() string {
	return "file"
}

type FileType int

const (
	FileTypeSystematic FileType = iota
	FileTypeProduct
	FileTypeBrand
	FileTypeAppPic
)

func FileTypeFromInt(value int) (FileType, error) {
	switch value {
	case int(FileTypeSystematic):
		return FileTypeSystematic, nil
	case int(FileTypeProduct):
		return FileTypeProduct, nil
	case int(FileTypeBrand):
		return FileTypeBrand, nil
	case int(FileTypeAppPic):
		return FileTypeAppPic, nil
	default:
		return 0, errors.New("invalid FileType value")
	}
}

func (ft FileType) GetFullInfo() (multiple, hasPriority, canForceDelete, isFileColumnNullable bool, table, mapTable, foreignColumn, fileColumn, priorityColumn, uploadDir, downloadPermission, listPermission, uploadPermission, deletePermission, changePriorirtyPermission string) {
	// default values
	fileColumn = "file_id"
	priorityColumn = "priority"
	multiple = false
	canForceDelete = true
	isFileColumnNullable = false

	switch ft {
	case FileTypeSystematic:
		downloadPermission = ACTION_FILE_SYSTEMATIC_DOWNLOAD
		uploadPermission = ACTION_FILE_SYSTEMATIC_UPLOAD
		listPermission = ACTION_FILE_SYSTEMATIC_LIST
		deletePermission = ACTION_FILE_SYSTEMATIC_DELETE

		multiple = true
		uploadDir = "uploads/systematic"

	case FileTypeProduct:
		downloadPermission = ACTION_FILE_PRODUCT_DOWNLOAD
		uploadPermission = ACTION_FILE_PRODUCT_UPLOAD
		listPermission = ACTION_FILE_PRODUCT_LIST
		deletePermission = ACTION_FILE_PRODUCT_DELETE
		changePriorirtyPermission = ACTION_FILE_PRODUCT_CHANGE_PRIORITY

		multiple = true
		hasPriority = true
		table = "product"
		mapTable = "product_file_map"
		foreignColumn = "product_id"
		uploadDir = "uploads/product"
		canForceDelete = false

	case FileTypeBrand:
		downloadPermission = ACTION_FILE_BRAND_DOWNLOAD
		uploadPermission = ACTION_FILE_BRAND_UPLOAD
		listPermission = ACTION_FILE_BRAND_LIST
		deletePermission = ACTION_FILE_BRAND_DELETE

		table = "brand"
		uploadDir = "uploads/brand"

	case FileTypeAppPic:
		downloadPermission = ACTION_FILE_APP_PIC_DOWNLOAD
		uploadPermission = ACTION_FILE_APP_PIC_UPLOAD
		listPermission = ACTION_FILE_APP_PIC_LIST
		deletePermission = ACTION_FILE_APP_PIC_DELETE

		table = "app_pic"
		uploadDir = "uploads/app-pic"
		hasPriority = true

	default:
		panic("invalid file type")
	}

	if !multiple {
		mapTable = table
	}

	return
}

func (ft FileType) GetInfo() (multiple, hasPriority bool, table, mapTable, foreignColumn, fileColumn, priorityColumn, uploadDir string) {
	multiple, hasPriority, _, _, table, mapTable, foreignColumn, fileColumn, priorityColumn, uploadDir, _, _, _, _, _ = ft.GetFullInfo()
	return
}

func (ft FileType) GetPermissions() (downloadPermission, listPermission, uploadPermission, deletePermission, changePriorirtyPermission string) {
	_, _, _, _, _, _, _, _, _, _, downloadPermission, listPermission, uploadPermission, deletePermission, changePriorirtyPermission = ft.GetFullInfo()
	return
}

func (ft FileType) GetUploadPermission() string {
	_, _, _, _, _, _, _, _, _, _, _, _, uploadPermission, _, _ := ft.GetFullInfo()
	return uploadPermission
}

func (ft FileType) GetDownloadPermission() string {
	_, _, _, _, _, _, _, _, _, _, downloadPermission, _, _, _, _ := ft.GetFullInfo()
	return downloadPermission
}

func (ft FileType) GetDeletePermission() string {
	_, _, _, _, _, _, _, _, _, _, _, _, _, deletePermission, _ := ft.GetFullInfo()
	return deletePermission
}

func (ft FileType) GetListPermission() string {
	_, _, _, _, _, _, _, _, _, _, _, listPermission, _, _, _ := ft.GetFullInfo()
	return listPermission
}

func (ft FileType) GetChangePriorityPermission() string {
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, changePriorirtyPermission := ft.GetFullInfo()
	return changePriorirtyPermission
}

func (ft FileType) GetDirectory() string {
	_, _, _, _, _, _, _, directory := ft.GetInfo()
	return directory
}

func (ft FileType) GetFileUrl(fileName string) string {
	return ft.GetDirectory() + "/" + fileName
}

func (ft FileType) GenerateWhereClause(db *gorm.DB, itemID uuid.UUID) *gorm.DB {
	multiple, hasPriority, table, mapTable, foreignColumn, fileColumn, priorityColumn, _ := ft.GetInfo()

	if multiple {
		db = db.Table(mapTable).Where(foreignColumn+"=?", itemID)
	} else {
		db = db.Table(table).Where("id=?", itemID)
	}

	if hasPriority {
		db = db.Order(priorityColumn + " ASC")
	}

	return db.Select(fileColumn)
}

func (ft FileType) CanForceDelete() bool {
	_, _, canForceDelete, _, _, _, _, _, _, _, _, _, _, _, _ := ft.GetFullInfo()
	return canForceDelete
}

func (ft FileType) IsFileColumnNullable() bool {
	_, _, _, isFileColumnNullable, _, _, _, _, _, _, _, _, _, _, _ := ft.GetFullInfo()
	return isFileColumnNullable
}
