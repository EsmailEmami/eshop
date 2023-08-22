package validations

import (
	"context"
	"errors"
	"reflect"

	dbpkg "github.com/esmailemami/eshop/db"
)

func ExistsInDB(model interface{}, column string, errorMsg string) func(value interface{}) error {
	return func(value interface{}) error {
		if value == nil {
			return nil
		}

		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return nil
			}
		}

		db := dbpkg.MustGormDBConn(context.Background())

		var count int64
		db.Model(model).
			Where(column+"=?", value).
			Count(&count)

		if count > 0 {
			return nil
		}

		return errors.New(errorMsg)
	}
}
func NotExistsInDB(model interface{}, column string, errorMsg string) func(value interface{}) error {
	return func(value interface{}) error {
		if value == nil {
			return nil
		}

		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return nil
			}
		}

		db := dbpkg.MustGormDBConn(context.Background())

		var count int64
		db.Model(model).
			Where(column+"=?", value).
			Count(&count)

		if count > 0 {
			return errors.New(errorMsg)
		}

		return nil
	}
}
