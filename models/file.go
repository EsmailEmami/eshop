package models

import (
	"errors"

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
	Products       []ProductFileMap `gorm:"foreignKey:file_id;references:id" json:"products"`
	Brands         []Brand          `gorm:"foreignKey:file_id;references:id" json:"brands"`
	AppPics        []AppPic         `gorm:"foreignKey:file_id;references:id" json:"appPics"`
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
func (ft FileType) GetInfo() (multiple, hasPriority bool, table, mapTable, foreignColumn, fileColumn, priorityColumn, uploadDir string) {
	// default values
	fileColumn = "file_id"
	priorityColumn = "priority"
	multiple = false

	switch ft {
	case FileTypeSystematic:
		multiple = true
		uploadDir = "uploads/systematic"

	case FileTypeProduct:
		multiple = true
		hasPriority = true
		table = "product"
		mapTable = "product_file_map"
		foreignColumn = "product_id"
		uploadDir = "uploads/product"

	case FileTypeBrand:
		table = "brand"
		uploadDir = "uploads/brand"

	case FileTypeAppPic:
		table = "app_pic"
		uploadDir = "uploads/app-pic"
		hasPriority = true

	default:
		panic("invalid file type")
	}

	if !multiple && hasPriority {
		mapTable = table
	}

	return
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
