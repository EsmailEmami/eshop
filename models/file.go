package models

import "errors"

type File struct {
	Model

	MimeType       string           `gorm:"mime_type"                        json:"mimeType"`
	Extension      string           `gorm:"extension"                        json:"extension"`
	OriginalName   string           `gorm:"original_name"                    json:"originalName"`
	UniqueFileName string           `gorm:"unique_file_name"                 json:"uniqueFineName"`
	FileType       FileType         `gorm:"file_type"                        json:"fileType"`
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
func (ft FileType) GetDirectory() string {
	switch ft {
	case FileTypeSystematic:
		return "uploads/systematic"
	case FileTypeProduct:
		return "uploads/product"
	case FileTypeBrand:
		return "uploads/brand"
	case FileTypeAppPic:
		return "uploads/app-pic"
	default:
		panic("invalid file type")
	}
}
