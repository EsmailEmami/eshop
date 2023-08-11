package middlewares

import (
	"context"
	"net/http"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/errors"
	"github.com/esmailemami/eshop/app/services/token"
	appDB "github.com/esmailemami/eshop/db"
	"github.com/google/uuid"
)

func Authentication(ctx *app.HttpContext) error {
	jwtToken, tokenString, err := token.LoadTokenFromHttpRequest(ctx.Request)
	if err != nil {
		return errors.NewUnauthorizedError(consts.UnauthorizedError, err)
	}

	// Get Token jit
	jit := (*jwtToken).JwtID()

	tokenID, err := uuid.Parse(jit)
	if err != nil {
		return errors.NewUnauthorizedError(consts.UnauthorizedError, err)
	}

	db := appDB.MustGormDBConn(ctx)

	// load token from DB
	authToken, user, err := token.GetValidAuthTokenByID(db, tokenID)
	if err != nil {
		return errors.NewUnauthorizedError(consts.UnauthorizedError, err)
	}

	// set contexts
	reqContext := context.WithValue(ctx.Request.Context(), consts.UserContext, *user)
	reqContext = context.WithValue(reqContext, consts.TokenContext, *jwtToken)
	reqContext = context.WithValue(reqContext, consts.TokenStringContext, tokenString)
	reqContext = context.WithValue(reqContext, consts.AuthTokenID, authToken.ID.String())

	ctx.Request = ctx.Request.WithContext(reqContext)
	return nil
}

func AuthenticationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewHttpContext(w, r)
		err := Authentication(ctx)
		if err != nil {
			app.ErrorResponseHandler(ctx, err)
			return
		}

		next.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	})
}
