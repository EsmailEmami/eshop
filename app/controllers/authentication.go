package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/models"
	"github.com/esmailemami/eshop/consts"
	dbpkg "github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/errors"
	dbmodels "github.com/esmailemami/eshop/models"
	"github.com/esmailemami/eshop/services/authentication"
	"github.com/esmailemami/eshop/services/notifier/email"
	"github.com/esmailemami/eshop/services/token"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
// @Router /admin/auth/login [post]
func LoginAdmin(ctx *app.HttpContext) error {
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
// @Router /user/auth/login [post]
func LoginUser(ctx *app.HttpContext) error {
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
// @Router /user/auth/is_authenticated [get]
func IsAuthenticated(ctx *app.HttpContext) error {
	userCtx, ok := ctx.Get(consts.UserContext)
	if !ok {
		return errors.NewBadRequestError(consts.UnauthorizedError, nil)
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
// @Router /user/auth/logout [get]
func Logout(ctx *app.HttpContext) error {
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

func logOutDone(ctx *app.HttpContext) error {
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
// @Router /user/auth/register [post]
func Register(ctx *app.HttpContext) error {
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

// Recovery Password godoc
// @Summary recovery user password.
// @Description recovery user password.
// @Tags Auth
// @Accept json
// @Produce json
// @Param RecoveryPasswordInput  body  models.RecoveryPasswordReqModel  true  "Recovery password model"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/auth/recoveryPasword [post]
func SendRecoveryPasswordRequest(ctx *app.HttpContext) error {
	var input models.RecoveryPasswordReqModel
	if err := ctx.BlindBind(&input); err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	if err := input.Validate(); err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	db := dbpkg.MustGormDBConn(ctx)

	var user dbmodels.User

	if err := db.Model(&dbmodels.User{}).
		Where("mobile=?", input.PhoneNumberOrEmailAddress).
		Or("email=?", input.PhoneNumberOrEmailAddress).First(&user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return errors.NewInternalServerError(consts.InternalServerError, err)
		} else {
			return ctx.QuickResponse(consts.RecoveryPasswordReqDone, http.StatusOK)
		}
	}

	if strings.TrimSpace(user.Email) == "" {
		return ctx.QuickResponse(consts.RecoveryPasswordReqDone, http.StatusOK)
	}

	// check that is there any not expired verification code
	if dbpkg.Exists(db, &dbmodels.VerificationCode{}, "scope=? AND key=? AND expire_at>? AND verified = false", dbmodels.VerificationCodeScopeEmail, user.Email, time.Now()) {
		return errors.NewValidationError("We have been sent you an email. If you do not receive the email try again later", nil)
	}

	verificaationCode := dbmodels.VerificationCode{
		BasicModel: dbmodels.BasicModel{
			ID: dbmodels.NewID(),
		},
		ExpireAt:   time.Now().Add(5 * time.Minute),
		MaxRetires: 3,
		Scope:      dbmodels.VerificationCodeScopeEmail,
		Key:        user.Email,
		Value:      uuid.NewString(),
		Verified:   false,
	}

	if err := db.Create(&verificaationCode).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	// send email
	notifier := email.NewNotifier("gmail")
	go func() {
		err := notifier.Send([]string{user.Email}, email.KeyForgotPassword, email.ForgotPassword{
			Username:    user.Username,
			RecoveryUrl: "http://127.0.0.1:3000/recoveryPassword/" + verificaationCode.Value,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	return ctx.QuickResponse(consts.RecoveryPasswordReqDone, http.StatusOK)
}

// Recovery Password godoc
// @Summary recovery user password.
// @Description recovery user password.
// @Tags Auth
// @Accept json
// @Produce json
// @Param key  path  string  true  "key"
// @Param RecoveryPasswordInput  body  models.RecoveryPasswordModel  true  "Recovery password model"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /user/auth/recoveryPasword/{key} [post]
func RecoveryPassword(ctx *app.HttpContext) error {
	key, err := uuid.Parse(ctx.GetPathParam("key"))
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	db := dbpkg.MustGormDBConn(ctx)

	var verificationCode dbmodels.VerificationCode

	if err := db.Model(&dbmodels.VerificationCode{}).Order("expire_at DESC").
		Find(&verificationCode, "value=? AND expire_at > now()", key).Error; err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	if verificationCode.Attempts > verificationCode.MaxRetires {
		return errors.NewBadRequestError("Your effort has exceeded the limit.", nil)
	}

	// save the user attempts
	verificationCode.Attempts++
	db.Save(&verificationCode)

	var input models.RecoveryPasswordModel
	if err := ctx.BlindBind(&input); err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}

	if err := input.Validate(); err != nil {
		return errors.NewValidationError(consts.ValidationError, err)
	}

	// encrypt password
	pass, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err := db.Model(&dbmodels.User{}).Where("email=?", verificationCode.Key).UpdateColumn("password", string(pass)).Error; err != nil {
		return errors.NewInternalServerError(consts.InternalServerError, err)
	}

	return ctx.QuickResponse(consts.OperationDone, http.StatusOK)
}
