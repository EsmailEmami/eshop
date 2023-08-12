package middlewares

import (
	"context"
	"net/http"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/errors"
	"github.com/esmailemami/eshop/app/services/authorization"
	"github.com/esmailemami/eshop/models"
)

// actions will check with OR condition.
// if one of them is allowed, it will return no error
func Permitted(actions ...string) func(ctx *app.HttpContext) error {
	return func(ctx *app.HttpContext) error {
		return authorization.CanAccess(ctx, actions...)
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
