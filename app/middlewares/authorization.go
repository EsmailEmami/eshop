package middlewares

import (
	"fmt"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/consts"
	"github.com/esmailemami/eshop/errors"
	"github.com/esmailemami/eshop/models"

	errs "errors"
)

// actions will check with OR condition.
// if one of them is allowed, it will return no error
func Permitted(actions ...string) func(ctx *app.HttpContext) error {
	return func(ctx *app.HttpContext) error {
		permitted := false

		// load token's related user
		user, ok := ctx.Request.Context().Value(consts.UserContext).(*models.User)
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
}
