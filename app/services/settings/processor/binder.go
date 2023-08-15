package processor

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

func Bind(db *sql.DB, model interface{}) error {
	var (
		rv = reflect.ValueOf(model).Elem()
		rt = reflect.TypeOf(rv.Interface())
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

	columnsName := getColumnsName(rt)

	numCols := rv.NumField()
	columns := make([]interface{}, numCols)
	for i := 0; i < numCols; i++ {
		field := rv.Field(i)
		columns[i] = field.Addr().Interface()
	}

	selectQry := fmt.Sprintf(`SELECT %s FROM "%s"."%s" LIMIT 1;`, strings.Join(columnsName, ","), schema, tableName)
	err = db.QueryRow(selectQry).Scan(columns...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	return nil
}
