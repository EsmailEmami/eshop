package authorization

import (
	"fmt"

	errs "errors"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/errors"
	"github.com/esmailemami/eshop/models"
)

func CanAccess(ctx *app.HttpContext, actions ...string) error {
	permitted := false

	// load token's related user
	user, ok := ctx.Request.Context().Value(consts.UserContext).(models.User)
	if !ok {
		return errors.NewUnauthorizedError(consts.UnauthorizedError, nil)
	}

	if user.Role == nil {
		return errors.NewForbiddenError(
			consts.ForbiddenError,
			errs.New("Permitted(): user.Role is nil"),
		)
	}

	// if one of the actions is permitted, so user should be able to continue
	for _, a := range actions {
		if user.Role.Permitted(a) {
			permitted = true
		}
	}

	if !permitted {
		return errors.NewForbiddenError(
			consts.ForbiddenError,
			fmt.Errorf(
				"Permitted(): user.Role is not permitted to do these actions: %v",
				actions,
			),
		)
	}
	return nil
}
