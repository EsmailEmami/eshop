package tablecreator

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
	gormSchema "gorm.io/gorm/schema"
)

type columnData struct {
	Type         string
	DefaultValue string
}

func CreateOrUpdate(db *sql.DB, model interface{}) error {
	rv := reflect.ValueOf(model)
	rt := reflect.TypeOf(rv.Elem().Interface())

	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("model must be a pointer")
	}

	tableName, schema, err := getTableData(rv)

	if err != nil {
		return err
	}

	defaultValuesMethod := reflect.ValueOf(model).MethodByName("LoadDefaultValues")
	if defaultValuesMethod.IsValid() {
		defaultValuesMethod.Call(nil)
	} else {
		return fmt.Errorf("the given struct must have TableName() function.")
	}

	exists, err := exists(db, schema, tableName)
	if err != nil {
		return err
	}

	columns := getColumnsData(rt, rv.Elem())
	existedColumns := make(map[string]string)

	if exists {
		existedColumns, err = getExistedColumn(db, tableName)
		if err != nil {
			return err
		}

		columnDefinitions := []string{}

		// add or alter columns
		for columnName, column := range columns {
			if _, ok := existedColumns[columnName]; ok {
				columnDefinitions = append(columnDefinitions, fmt.Sprintf("DROP COLUMN IF EXISTS %s", columnName))
			}
			columnDefinitions = append(columnDefinitions, fmt.Sprintf("ADD %s %s", columnName, column.Type))
		}

		// delete columns
		for columnName := range existedColumns {
			if _, ok := columns[columnName]; !ok {
				columnDefinitions = append(columnDefinitions, fmt.Sprintf("DROP COLUMN IF EXISTS %s", columnName))
			}
		}

		sqlCommand := fmt.Sprintf("ALTER TABLE %s.%s %s;", schema, tableName, strings.Join(columnDefinitions, ","))

		_, err = db.Exec(sqlCommand)
		if err != nil {
			return err
		}
	} else {
		columnDefinitions := []string{}

		// add columns
		for columnName, column := range columns {
			columnDefinitions = append(columnDefinitions, fmt.Sprintf("%s %s", columnName, column.Type))
		}

		sqlCommand := fmt.Sprintf(`CREATE TABLE "%s"."%s" (%s);`, schema, tableName, strings.Join(columnDefinitions, ","))
		_, err = db.Exec(sqlCommand)

		if err != nil {
			return err
		}
	}

	if !exists {
		columnsName := []string{}
		values := []string{}

		for columnName, column := range columns {
			columnsName = append(columnsName, columnName)
			values = append(values, column.DefaultValue)
		}

		insertComand := fmt.Sprintf(`INSERT INTO "%s"."%s" (%s) VALUES(%s)`, schema, tableName, strings.Join(columnsName, ","), strings.Join(values, ","))
		db.Exec(insertComand)
	} else {
		columnsUpdate := []string{}

		for columnName, column := range columns {
			if _, ok := existedColumns[columnName]; !ok {
				columnsUpdate = append(columnsUpdate, fmt.Sprintf(`"%s" = N"%s"`, columnName, column.DefaultValue))
			}
		}

		updateCommand := fmt.Sprintf(`UPDATE "%s"."%s" SET %s;`, schema, tableName, strings.Join(columnsUpdate, ","))
		db.Exec(updateCommand)
	}

	return nil
}

func exists(db *sql.DB, schema, tableName string) (bool, error) {
	row := db.QueryRow("SELECT EXISTS (SELECT 1 FROM pg_tables WHERE schemaname = $1 AND tablename = $2);", schema, tableName)

	if row.Err() != nil {
		return false, row.Err()
	}

	var exists bool

	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func getTableData(rv reflect.Value) (tableName, schema string, err error) {
	tableNameMethod := rv.MethodByName("TableName")
	if tableNameMethod.IsValid() {
		tableName = tableNameMethod.Call(nil)[0].String()
	} else {
		return "", "", fmt.Errorf("the given struct must have TableName() function.")
	}

	schemaMethod := rv.MethodByName("SchemaName")
	if schemaMethod.IsValid() {
		schema = schemaMethod.Call(nil)[0].String()
	} else {
		return "", "", fmt.Errorf("the given struct must have SchemaName() function.")
	}

	return
}

func getExistedColumn(db *sql.DB, tableName string) (map[string]string, error) {
	rows, err := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = $1", tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make(map[string]string)

	for rows.Next() {
		var columnName, dataType string
		if err := rows.Scan(&columnName, &dataType); err != nil {
			return nil, err
		}
		data[columnName] = dataType
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func getColumnsData(rt reflect.Type, rv reflect.Value) map[string]*columnData {
	data := make(map[string]*columnData)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		columnType := getColumnType(field.Type)
		dbInfo := gormSchema.ParseTagSetting(field.Tag.Get("gorm"), ";")
		columnName, columnOk := dbInfo["COLUMN"]

		if !columnOk {
			columnName = field.Name
		}

		columnData := &columnData{
			Type: columnType,
		}

		fieldValue := rv.Field(i)

		val, err := interfaceToString(fieldValue.Interface())

		if err == nil {
			columnData.DefaultValue = val
		}

		data[columnName] = columnData
	}

	return data
}

func getColumns(rt reflect.Type) []string {
	data := []string{}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		dbInfo := gormSchema.ParseTagSetting(field.Tag.Get("gorm"), ";")
		columnName, columnOk := dbInfo["COLUMN"]

		if !columnOk {
			columnName = field.Name
		}

		data = append(data, columnName)
	}

	return data
}

func getColumnType(fieldType reflect.Type) string {
	switch fieldType.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "INT"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "INT"
	case reflect.Float32, reflect.Float64:
		return "NUMERIC"
	case reflect.String:
		return "VARCHAR(512)"
	case reflect.Ptr:
		return getColumnType(fieldType.Elem())
	case reflect.Struct:
		if fieldType == reflect.TypeOf(uuid.UUID{}) {
			return "UUID"
		}
	default:
		return "VARCHAR(512)"
	}

	return "VARCHAR(512)"
}

func interfaceToString(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case int8:
		return strconv.FormatInt(int64(v), 10), nil
	case int16:
		return strconv.FormatInt(int64(v), 10), nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		return interfaceToString(rv.Elem().Interface())
	}

	if rv.Kind() == reflect.Struct || (rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Struct) {
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	}

	return "NULL", nil
}
