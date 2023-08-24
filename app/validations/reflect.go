package validations

import "reflect"

func IsNil(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return true
		}
	}

	return false
}

func Value(value interface{}) interface{} {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		} else {
			return v.Elem().Interface()
		}
	}
	return value
}
