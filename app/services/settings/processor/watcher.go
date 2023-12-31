package processor

import "reflect"

type SettingItem struct {
	Field       string  `json:"field"`
	FieldType   string  `json:"fieldType"`
	Value       any     `json:"value"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	IsNullable  bool    `json:"isNullable"`
}

func GetItems(model any) []SettingItem {
	var (
		rv    = reflect.ValueOf(model).Elem()
		rt    = reflect.TypeOf(rv.Interface())
		items = []SettingItem{}
	)

	for i := 0; i < rt.NumField(); i++ {
		var (
			field          = rt.Field(i)
			fieldVal       = rv.Field(i)
			value      any = nil
			isNullable     = false
			fieldType  string
		)

		if fieldVal.Kind() == reflect.Ptr {
			fieldType = field.Type.Elem().String()
			isNullable = true
			if !fieldVal.IsNil() {
				value = fieldVal.Elem().Interface()
			}
		} else {
			value = fieldVal.Interface()
			fieldType = field.Type.String()
		}

		item := SettingItem{
			Field:       field.Name,
			Value:       value,
			Title:       getColumnTitle(field),
			Description: getColumnDescription(field),
			IsNullable:  isNullable,
			FieldType:   fieldType,
		}

		items = append(items, item)
	}

	return items
}

func getColumnTitle(field reflect.StructField) *string {
	columnTitle, columnOk := field.Tag.Lookup("title")
	if !columnOk {
		return nil
	}
	return &columnTitle
}

func getColumnDescription(field reflect.StructField) *string {
	columnDescription, columnOk := field.Tag.Lookup("description")
	if !columnOk {
		return nil
	}
	return &columnDescription
}
