package validations

import (
	"errors"
	"regexp"

	"github.com/esmailemami/eshop/app/consts"
)

func NumericValidator() func(value interface{}) error {
	return func(value interface{}) error {
		if value == nil {
			return nil
		}
		v, ok := value.(string)
		if !ok {
			return errors.New(consts.InvalidNumeric)
		}

		if v == "" {
			return nil
		}

		matched, _ := regexp.Match("^[0-9]+$", []byte(v))
		if !matched {
			return errors.New(consts.InvalidNumeric)
		}

		return nil
	}
}
