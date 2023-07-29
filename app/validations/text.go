package validations

import (
	"errors"
	"strings"

	"github.com/esmailemami/eshop/consts"
	"github.com/esmailemami/eshop/services/sanitize"
)

// validates that the given input is a clear text or not
// it allows a-z, A-Z, 0-9, persian and arabic letters
func ClearText() func(value interface{}) error {
	return func(value interface{}) error {
		if value == nil {
			return nil
		}
		txt, ok := value.(string)
		if !ok {
			return errors.New(consts.InvalidCharacters)
		}
		if sanitize.AsClearText(txt) == strings.Trim(txt, " ") {
			return nil
		}
		return errors.New(consts.InvalidCharacters)
	}
}

// validates that the given input is a clear text or not
// Username allowed characters: 0-9 a-z(all lower case)
func UserName() func(value interface{}) error {
	return func(value interface{}) error {
		if value == nil {
			return nil
		}
		txt, ok := value.(string)
		if !ok {
			return errors.New(consts.InvalidCharacters)
		}
		if sanitize.AsUsername(txt) == strings.Trim(txt, " ") {
			return nil
		}
		return errors.New(consts.InvalidCharacters)
	}
}

// validate that the given input can be set as Code
func Code() func(value interface{}) error {
	return func(value interface{}) error {
		if value == nil {
			return nil
		}
		txt, ok := value.(string)
		if !ok {
			return errors.New(consts.InvalidCharacters)
		}
		if sanitize.AsCode(txt) == strings.Trim(txt, " ") {
			return nil
		}
		return errors.New(consts.InvalidCharacters)
	}
}
