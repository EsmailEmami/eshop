package controllers

import (
	"net/http"

	"github.com/esmailemami/eshop/apphttp"
	"github.com/esmailemami/eshop/apphttp/models"
	"github.com/esmailemami/eshop/consts"
	dbpkg "github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/errors"
	dbmodels "github.com/esmailemami/eshop/models"
	"github.com/esmailemami/eshop/services/authentication"
	"github.com/esmailemami/eshop/services/token"
	"github.com/google/uuid"
)

// LoginCustomer godoc
// @Summary Log user in
// @Description Log user in
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginInput   body  models.LoginInputModel  true  "Login input model"
// @Success 200 {object} models.LoginOutputModel
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /auth/login [post]
func Login(ctx *apphttp.HttpContext) error {
	var input models.LoginInputModel

	if err := ctx.BlindBind(&input); err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	if err := input.Validate(); err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	input.IP = ctx.ClientIP()
	input.UserAgent = ctx.UserAgent()

	output, err := authentication.LoginByUsername(ctx, input)
	if err != nil {
		return errors.NewBadRequestError(err.Error(), err)
	}

	ctx.SetCookie("Authorization", output.Token, int(output.ExpiresIn), "/", "", true, true)

	return ctx.JSON(*output, http.StatusOK)
}

// IsAuthenticated godoc
// @Summary Returns the logged in user
// @Description Returns the logged in user.
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} []dbmodels.User
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /auth/is_authenticated [get]
func IsAuthenticated(ctx *apphttp.HttpContext) error {
	userCtx, ok := ctx.Get(consts.UserContext)
	if !ok {
		return errors.NewBadRequestError("کاربر یافت نشد", nil)
	}

	user := userCtx.(dbmodels.User)
	if user.RoleID != nil {
		// user role data not loaded
		if user.Role == nil {
			dbConn := dbpkg.MustGormDBConn(ctx)
			role := dbmodels.Role{}
			dbConn.Where("id", user.RoleID).First(&role)
			user.Role = &role
		}
	}
	return ctx.JSON(user, http.StatusOK)
}

// Logout godoc
// @Summary Logout the user.
// @Description Logout the user.
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /auth/logout [get]
func Logout(ctx *apphttp.HttpContext) error {
	jwtToken, _, err := token.LoadTokenFromHttpRequest(ctx.Request)
	if err != nil {
		return logOutDone(ctx)
	}

	jit := uuid.MustParse((*jwtToken).JwtID())

	// load token from DB
	err = authentication.RevokeAuthTokenByID(dbpkg.MustGormDBConn(ctx), jit)
	if err != nil {
		// can't revoke token, so let client do logout process
		// and avoid revoking token from DB
		return logOutDone(ctx)
	}

	return logOutDone(ctx)
}

func logOutDone(ctx *apphttp.HttpContext) error {
	ctx.SetCookie("Authorization", "", 0, "/", "", true, true)
	return ctx.QuickResponse(consts.LoggedOut, http.StatusOK)
}

// Register godoc
// @Summary a new user.
// @Description register a new user.
// @Tags Auth
// @Accept json
// @Produce json
// @Param RegisterInput  body  models.RegisterInputModel  true  "Register input model"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /auth/register [post]
func Register(ctx *apphttp.HttpContext) error {
	var input models.RegisterInputModel
	if err := ctx.BlindBind(&input); err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	if err := input.Validate(); err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	user := input.ToDBModel()

	db := dbpkg.MustGormDBConn(ctx)

	if err := db.Create(user).Error; err != nil {
		return errors.NewValidationError(consts.InternalServerError, err)
	}

	return ctx.QuickResponse(consts.RegistrationDone, http.StatusOK)
}
