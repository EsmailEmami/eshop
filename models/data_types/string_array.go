package datatypes

import (
	"database/sql/driver"
	"encoding/json"
)

type StringArray []string

func (p StringArray) Value() (driver.Value, error) {
	valueString, err := json.Marshal(p)
	return string(valueString), err
}

func (j *StringArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	var bts []byte
	switch v := value.(type) {
	case []byte:
		bts = v
	case string:
		bts = []byte(v)
	case nil:
		*j = nil
		return nil
	}
	return json.Unmarshal(bts, &j)
}
