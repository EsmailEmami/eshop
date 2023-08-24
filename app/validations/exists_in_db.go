package validations

import (
	"context"
	"errors"
	"reflect"

	dbpkg "github.com/esmailemami/eshop/db"
	"github.com/google/uuid"
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

func NotExistsInDBWithID(model interface{}, column string, id uuid.UUID, errorMsg string) func(value interface{}) error {
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
			Where("id != ", id).
			Count(&count)

		if count > 0 {
			return errors.New(errorMsg)
		}

		return nil
	}
}

func NotExistsInDBWithCond(model interface{}, column string, errorMsg string, condition interface{}, args ...interface{}) func(value interface{}) error {
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
			Where(condition, args...).
			Count(&count)

		if count > 0 {
			return errors.New(errorMsg)
		}

		return nil
	}
}

func ExistsInDBWithCond(model interface{}, column string, errorMsg string, condition interface{}, args ...interface{}) func(value interface{}) error {
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
			Where(condition, args...).
			Count(&count)

		if count > 0 {
			return nil
		}

		return errors.New(errorMsg)
	}
}
