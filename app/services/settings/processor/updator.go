package processor

import (
	"database/sql"
	"fmt"
	"reflect"
)

func UpdateField(db *sql.DB, model interface{}, field string, value interface{}) error {
	var (
		rv = reflect.ValueOf(model).Elem()
		rt = rv.Type()
	)

	tableName, schema, err := getTableData(rv)
	if err != nil {
		return err
	}

	tableExists, err := tableExists(db, schema, tableName)
	if err != nil {
		return err
	}
	if !tableExists {
		return fmt.Errorf("cannot find table %s", tableName)
	}

	fieldType, found := rt.FieldByName(field)
	if !found {
		return fmt.Errorf("cannot find field %s", field)
	}
	column := getColumnName(fieldType)
	_ = column

	fieldV := rv.FieldByName(field)
	if !fieldV.CanSet() {
		return fmt.Errorf("cannot set value for field %s", field)
	}
	if fieldValue, err := convertToFieldType(value, fieldType.Type); err == nil {
		fieldV.Set(reflect.ValueOf(fieldValue))
	} else {
		return err
	}

	dbStringVal := interfaceToString(value)

	updateCommand := fmt.Sprintf(
		`UPDATE "%s"."%s" SET "%s" = %s;`,
		schema,
		tableName,
		column,
		dbStringVal,
	)
	_, err = db.Exec(updateCommand)

	return err
}

func convertToFieldType(value interface{}, targetType reflect.Type) (interface{}, error) {
	targetValue := reflect.New(targetType).Elem()
	sourceValue := reflect.ValueOf(value)

	if targetType.Kind() == reflect.Ptr {
		if !sourceValue.IsValid() {
			targetValue.Set(reflect.Zero(targetType))
			return targetValue.Interface(), nil
		} else if sourceValue.Type().ConvertibleTo(targetType.Elem()) {
			convertedValue := reflect.New(targetType.Elem()).Elem()
			convertedValue.Set(sourceValue.Convert(targetType.Elem()))
			targetValue.Set(convertedValue.Addr())
			return targetValue.Interface(), nil
		}
	} else if sourceValue.IsValid() && sourceValue.Type().ConvertibleTo(targetType) {
		targetValue.Set(sourceValue.Convert(targetType))
		return targetValue.Interface(), nil
	}

	return nil, fmt.Errorf("cannot convert value to target type")
}
