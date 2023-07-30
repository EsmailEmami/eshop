package models

import "github.com/google/uuid"

type AppPic struct {
	Model

	Priority    int        `gorm:"priority"                         json:"priority"`
	AppPicType  AppPicType `gorm:"app_pic_type"                     json:"appPicType"`
	FileID      uuid.UUID  `gorm:"file_id"                          json:"fileId"`
	File        *File      `gorm:"foreignKey:file_id;references:id" json:"file"`
	Title       string     `gorm:"title"                            json:"title"`
	Description string     `gorm:"description"                      json:"description"`
}

func (AppPic) TableName() string {
	return "app_pic"
}

type AppPicType int

const (
	AppPicTypeSlider AppPicType = iota
)
