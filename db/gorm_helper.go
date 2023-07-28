package db

import "gorm.io/gorm"

func Paginate(db *gorm.DB, page, size int) *gorm.DB {
	return db.Offset(size * (page - 1)).Limit(size)
}

func Exists(db *gorm.DB, model interface{}, condition interface{}, args ...interface{}) bool {
	var count int64
	db.Model(model).Where(condition, args...).Count(&count)
	return count > 0
}
