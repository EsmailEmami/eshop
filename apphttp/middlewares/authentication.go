package middlewares

import (
	"context"
	"net/http"

	"github.com/esmailemami/eshop/apphttp"
	"github.com/esmailemami/eshop/consts"
	appDB "github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/errors"
	"github.com/esmailemami/eshop/services/token"
	"github.com/google/uuid"
)

func Authentication(ctx *apphttp.HttpContext) error {
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
		ctx := apphttp.NewHttpContext(w, r)
		err := Authentication(ctx)
		if err != nil {
			apphttp.ErrorResponseHandler(ctx, err)
			return
		}

		next.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	})
}
