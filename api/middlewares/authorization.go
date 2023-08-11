package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/errors"
	"github.com/esmailemami/eshop/models"

	errs "errors"
)

// actions will check with OR condition.
// if one of them is allowed, it will return no error
func Permitted(actions ...string) func(ctx *app.HttpContext) error {
	return func(ctx *app.HttpContext) error {
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
}

func throwForbiddenError(ctx *app.HttpContext) {
	app.ErrorResponseHandler(ctx, errors.NewForbiddenError(consts.ForbiddenError, nil))
}

// only users with administrative roles can call these apis
func CanInvokeRouteUnlessAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewHttpContext(w, r)

		userCtx, ok := ctx.Get(consts.UserContext)
		if !ok {
			throwForbiddenError(ctx)
			return
		}

		user := userCtx.(models.User)

		if user.Role == nil {
			throwForbiddenError(ctx)
			return
		}

		if !user.Role.Permitted(models.ACTION_CAN_LOGIN_ADMIN) {
			throwForbiddenError(ctx)
		}

		reqContext := context.WithValue(ctx.Request.Context(), consts.UserActAsContext, consts.UserActAsAdmin)
		ctx.Request = ctx.Request.WithContext(reqContext)

		next.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	})
}

// only users with administrative roles can call these apis
func CanInvokeRouteUnlessUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewHttpContext(w, r)

		userCtx, ok := ctx.Get(consts.UserContext)
		if !ok {
			throwForbiddenError(ctx)
			return
		}

		user := userCtx.(models.User)

		if user.Role == nil {
			throwForbiddenError(ctx)
			return
		}

		if !user.Role.Permitted(models.ACTION_CAN_LOGIN_USER) {
			throwForbiddenError(ctx)
		}

		reqContext := context.WithValue(ctx.Request.Context(), consts.UserActAsContext, consts.UserActAsUser)
		ctx.Request = ctx.Request.WithContext(reqContext)

		next.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	})
}
