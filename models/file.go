package models

type File struct {
	Model

	MimeType       string           `gorm:"mime_type"                        json:"mimeType"`
	OriginalName   string           `gorm:"original_name"                    json:"originalName"`
	UniqueFileName string           `gorm:"unique_file_name"                 json:"uniqueFineName"`
	FileType       FileType         `gorm:"file_type"                        json:"fileType"`
	Products       []ProductFileMap `gorm:"foreignKey:file_id;references:id" json:"products"`
}

type FileType int

const (
	FileTypeProduct FileType = iota
)
