package models

import "github.com/google/uuid"

type AppPic struct {
	Model

	Priority    int        `gorm:"column:priority"                  json:"priority"`
	AppPicType  AppPicType `gorm:"column:app_pic_type"              json:"appPicType"`
	FileID      uuid.UUID  `gorm:"column:file_id"                   json:"fileId"`
	File        *File      `gorm:"foreignKey:file_id;references:id" json:"file"`
	Title       string     `gorm:"column:title"                     json:"title"`
	Description string     `gorm:"column:description"               json:"description"`
	Url         string     `gorm:"url"                              json:"url"`
}

func (AppPic) TableName() string {
	return "app_pic"
}

type AppPicType int

const (
	AppPicTypeSlider AppPicType = iota
	AppPicTypeSection
	AppPicTypeBilboard
)
