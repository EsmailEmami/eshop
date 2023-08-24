package validations

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/esmailemami/eshop/app/consts"
)

func IsValidNationalCode() func(value interface{}) error {
	return func(value interface{}) error {
		if IsNil(value) {
			return nil
		}

		nationalCode, ok := Value(value).(string)
		if !ok {
			return errors.New(consts.InvalidNationalCode)
		}

		// Check if the input is exactly 10 digits
		if len(nationalCode) != 10 {
			return errors.New(consts.InvalidNationalCode)
		}

		// Check if the input consists of digits only
		match, _ := regexp.MatchString("^[0-9]+$", nationalCode)
		if !match {
			return errors.New(consts.InvalidNationalCode)
		}

		// Calculate the checksum digit
		checkDigit, _ := strconv.Atoi(nationalCode[9:10])
		coefficients := []int{10, 9, 8, 7, 6, 5, 4, 3, 2}
		sum := 0

		for i := 0; i < 9; i++ {
			digit, _ := strconv.Atoi(nationalCode[i : i+1])
			sum += digit * coefficients[i]
		}

		remainder := sum % 11
		if remainder < 2 && checkDigit == remainder {
			return nil
		} else if remainder >= 2 && checkDigit == 11-remainder {
			return nil
		}

		return errors.New(consts.InvalidNationalCode)
	}
}
