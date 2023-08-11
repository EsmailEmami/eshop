package validations

import (
	"errors"
	"strings"

	"github.com/esmailemami/eshop/app/consts"
)

func TimeValidator(value interface{}) error {
	str := value.(string)
	if len(str) != 5 {
		return errors.New(str + ":" + consts.InvalidTime)
	}
	parts := strings.Split(str, ":")
	if len(parts) != 2 {
		return errors.New(str + ":" + consts.InvalidTime)
	}

	if parts[0] > "23" || parts[0] < "00" {
		return errors.New(str + ":" + consts.InvalidTime)
	}

	if parts[1] > "59" || parts[1] < "00" {
		return errors.New(str + ":" + consts.InvalidTime)
	}

	return nil
}
