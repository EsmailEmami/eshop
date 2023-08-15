package processor

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const defaultMaximumLength = 512

type columnInfo struct {
	Type         string
	DefaultValue string
}

func CreateOrUpdate(db *sql.DB, model interface{}) error {
	var (
		rv = reflect.ValueOf(model)
		rt = reflect.TypeOf(rv.Elem().Interface())
	)

	if rv.Kind() != reflect.Ptr {
		return errors.New("model must be a pointer")
	}

	tableName, schema, err := getTableData(rv)

	if err != nil {
		return err
	}

	if err := callDefaultValuesMethod(rv); err != nil {
		return err
	}

	tableExists, err := tableExists(db, schema, tableName)
	if err != nil {
		return err
	}

	var (
		columns        = getColumnsInfo(rt, rv.Elem())
		existedColumns = make(map[string]string)
		sqlCommand     string
	)

	if tableExists {
		existedColumns, err = getExistedColumn(db, schema, tableName)
		if err != nil {
			return err
		}

		columnDefinitions := []string{}

		// add or alter columns
		for columnName, column := range columns {
			if existedColumnType, ok := existedColumns[columnName]; ok {
				if existedColumnType != column.Type {
					columnDefinitions = append(columnDefinitions, fmt.Sprintf(`DROP COLUMN IF EXISTS "%s"`, columnName))
					columnDefinitions = append(columnDefinitions, fmt.Sprintf(`ADD "%s" %s`, columnName, column.Type))
				}
			} else {
				columnDefinitions = append(columnDefinitions, fmt.Sprintf(`ADD "%s" %s`, columnName, column.Type))
			}
		}

		// delete columns that are not existed
		for columnName := range existedColumns {
			if _, ok := columns[columnName]; !ok {
				columnDefinitions = append(columnDefinitions, fmt.Sprintf(`DROP COLUMN IF EXISTS "%s"`, columnName))
			}
		}

		if len(columnDefinitions) > 0 {
			sqlCommand = fmt.Sprintf(`ALTER TABLE "%s"."%s" %s;`, schema, tableName, strings.Join(columnDefinitions, ","))
		}

	} else {
		columnDefinitions := []string{}

		// add columns
		for columnName, column := range columns {
			columnDefinitions = append(columnDefinitions, fmt.Sprintf(`"%s" %s`, columnName, column.Type))
		}

		sqlCommand = fmt.Sprintf(`CREATE TABLE "%s"."%s" (%s);`, schema, tableName, strings.Join(columnDefinitions, ","))
	}

	// create or update table
	if strings.TrimSpace(sqlCommand) != "" {
		_, err = db.Exec(sqlCommand)
		if err != nil {
			return err
		}
	}

	dataCommand := ""

	if !tableExists {
		var (
			columnsName = []string{}
			values      = []string{}
		)

		for columnName, column := range columns {
			columnsName = append(columnsName, fmt.Sprintf(`"%s"`, columnName))
			values = append(values, column.DefaultValue)
		}

		dataCommand = fmt.Sprintf(`INSERT INTO "%s"."%s" (%s) VALUES(%s)`, schema, tableName, strings.Join(columnsName, ","), strings.Join(values, ","))
	} else {
		columnsUpdate := []string{}

		for columnName, column := range columns {
			if _, ok := existedColumns[columnName]; !ok {
				columnsUpdate = append(columnsUpdate, fmt.Sprintf(`"%s" = %s`, columnName, column.DefaultValue))
			}
		}

		if len(columnsUpdate) > 0 {
			dataCommand = fmt.Sprintf(`UPDATE "%s"."%s" SET %s;`, schema, tableName, strings.Join(columnsUpdate, ","))
		}
	}

	if strings.TrimSpace(dataCommand) != "" {
		_, err := db.Exec(dataCommand)
		if err != nil {
			return err
		}
	}

	return nil
}

func callDefaultValuesMethod(rv reflect.Value) error {
	defaultValuesMethod := rv.MethodByName("LoadDefaultValues")
	if defaultValuesMethod.IsValid() {
		defaultValuesMethod.Call(nil)
	}

	return nil
}

func tableExists(db *sql.DB, schema, tableName string) (bool, error) {
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
		return "", "", fmt.Errorf("the given struct must have TableName() function")
	}

	schemaMethod := rv.MethodByName("SchemaName")
	if schemaMethod.IsValid() {
		schema = schemaMethod.Call(nil)[0].String()
	} else {
		return "", "", fmt.Errorf("the given struct must have SchemaName() function")
	}

	return
}

func getExistedColumn(db *sql.DB, schema, tableName string) (map[string]string, error) {
	rows, err := db.Query("SELECT column_name, data_type,character_maximum_length FROM information_schema.columns WHERE table_schema = $1 AND table_name = $2", schema, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make(map[string]string)

	for rows.Next() {
		var (
			columnName, dataType string
			maximumLength        sql.NullInt32
		)
		if err := rows.Scan(&columnName, &dataType, &maximumLength); err != nil {
			return nil, err
		}
		data[columnName] = castColumnType(dataType, int(maximumLength.Int32))
	}

	return data, nil
}

func getColumnsInfo(rt reflect.Type, rv reflect.Value) map[string]*columnInfo {
	data := make(map[string]*columnInfo)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		columnType := getColumnType(field.Type, getMaximumLength(field))
		columnName := getColumnName(field)

		columnInfo := &columnInfo{
			Type:         columnType,
			DefaultValue: interfaceToString(rv.Field(i).Interface()),
		}

		data[columnName] = columnInfo
	}

	return data
}

func getColumnsName(rt reflect.Type) []string {
	data := []string{}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		data = append(data, getColumnName(field))
	}

	return data
}

func getColumnName(field reflect.StructField) string {
	columnName, columnOk := field.Tag.Lookup("column")
	if !columnOk {
		columnName = field.Name
	}

	return columnName
}

func getColumnType(fieldType reflect.Type, maximumLength int) string {
	switch fieldType.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "INT"
	case reflect.Float32, reflect.Float64:
		return "NUMERIC"
	case reflect.String:
		return fmt.Sprintf("VARCHAR(%d)", maximumLength)
	case reflect.Ptr:
		return getColumnType(fieldType.Elem(), maximumLength)
	case reflect.Struct:
		if fieldType == reflect.TypeOf(uuid.UUID{}) {
			return "UUID"
		}
	}

	panic("Invalid column type")
}

func getMaximumLength(field reflect.StructField) int {
	columnLength, columnOk := field.Tag.Lookup("maximumLength")
	if !columnOk {
		columnLength = strconv.Itoa(defaultMaximumLength)
	}

	intVal, err := strconv.Atoi(columnLength)
	if err != nil {
		return defaultMaximumLength
	}
	return intVal
}

func interfaceToString(value interface{}) string {
	if value == nil {
		return "NULL"
	}

	val := ""
	switch v := value.(type) {
	case string:
		val = "'" + v + "'"
	case int:
		val = strconv.Itoa(v)
	case int8:
	case int16:
	case int32:
		val = strconv.FormatInt(int64(v), 10)
	case int64:
		val = strconv.FormatInt(v, 10)
	case uint:
	case uint8:
	case uint16:
	case uint32:
		val = strconv.FormatUint(uint64(v), 10)
	case uint64:
		val = strconv.FormatUint(v, 10)
	case float32:
		val = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		val = strconv.FormatFloat(v, 'f', -1, 64)
	}

	if val == "" {
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Ptr {
			if rv.Elem().IsValid() {
				val = interfaceToString(rv.Elem().Interface())
			}

			if val == "''" || val == "" {
				val = "NULL"
			}
		}

		if rv.Kind() == reflect.Struct || (rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Struct) {
			jsonBytes, err := json.Marshal(value)
			if err != nil {
				val = ""
			}
			val = string(jsonBytes)
		}
	}

	return val
}

func castColumnType(columnType string, maximumLength int) string {
	switch strings.ToLower(columnType) {
	case "character varying":
		if maximumLength == 0 {
			return "VARCHAR"
		} else {
			return fmt.Sprintf("VARCHAR(%d)", maximumLength)
		}

	case "integer":
		return "INT"

	case "text":
		return "TEXT"

	case "timestamp":
		return "TIMESTAMP"

	case "boolean":
		return "BOOLEAN"

	case "numeric":
		return "NUMERIC"

	case "json":
		return "JSON"

	case "jsonb":
		return "JSONB"

	case "uuid":
		return "UUID"

	case "bytea":
		return "BYTEA"

	case "smallint":
		return "SMALLINT"

	case "bigint":
		return "BIGINT"

	case "real":
		return "REAL"

	case "double precision":
		return "DOUBLE PRECISION"

	case "date":
		return "DATE"

	case "time":
		return "TIME"

	case "interval":
		return "INTERVAL"

	case "serial":
		return "SERIAL"

	case "bigserial":
		return "BIGSERIAL"

	case "money":
		return "MONEY"

	case "point":
		return "POINT"

	case "line":
		return "LINE"

	case "lseg":
		return "LSEG"

	case "box":
		return "BOX"

	case "path":
		return "PATH"

	case "polygon":
		return "POLYGON"

	case "circle":
		return "CIRCLE"

	case "inet":
		return "INET"

	case "cidr":
		return "CIDR"

	case "macaddr":
		return "MACADDR"

	default:
		panic("Invalid column type")
	}
}
