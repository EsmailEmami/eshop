package validations

import (
	"errors"
	"regexp"

	"github.com/esmailemami/eshop/app/consts"
)

func IsValidPostalCode() func(value interface{}) error {
	return func(value interface{}) error {
		postalCode, ok := value.(string)
		if !ok {
			return errors.New(consts.InvalidPostalCode)
		}

		// Check if the input is exactly 10 digits
		if len(postalCode) != 10 {
			return errors.New(consts.InvalidPostalCode)
		}

		// Check if the input consists of digits only
		match, _ := regexp.MatchString("^[0-9]+$", postalCode)
		if !match {
			return errors.New(consts.InvalidPostalCode)
		}

		return nil
	}
}
