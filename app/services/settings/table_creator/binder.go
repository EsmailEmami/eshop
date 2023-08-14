package tablecreator

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

func Bind[T any](db *sql.DB, model T) (*T, error) {
	rv := reflect.ValueOf(model)
	rt := reflect.TypeOf(model)

	tableName, schema, err := getTableData(rv)

	if err != nil {
		return nil, err
	}

	exists, err := exists(db, schema, tableName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("cannot find table %s", tableName)
	}

	columnsName := getColumns(rt)
	s := reflect.ValueOf(&model).Elem()
	numCols := s.NumField()
	columns := make([]interface{}, numCols)
	for i := 0; i < numCols; i++ {
		field := s.Field(i)
		columns[i] = field.Addr().Interface()
	}

	selectQry := fmt.Sprintf(`SELECT %s FROM "%s"."%s" LIMIT 1;`, strings.Join(columnsName, ","), schema, tableName)
	err = db.QueryRow(selectQry).Scan(columns...)

	if err != nil {
		if err == sql.ErrNoRows {
			return &model, nil
		}
		return nil, err
	}

	return &model, nil
}
