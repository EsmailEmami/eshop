package authentication

import (
	"context"
	"errors"
	"time"

	appmodels "github.com/esmailemami/eshop/app/models"
	"github.com/esmailemami/eshop/consts"
	dbpkg "github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/models"
	"github.com/esmailemami/eshop/services/token"
	"github.com/esmailemami/eshop/services/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginByUsername(ctx context.Context, input appmodels.LoginInputModel) (*appmodels.LoginOutputModel, error) {
	db := dbpkg.MustGormDBConn(ctx)

	user, err := user.GetUserByUsername(input.Username)

	if err != nil {
		return nil, errors.New(consts.LoginFailed)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, errors.New(consts.LoginFailed)
	}

	tx := db.Begin()
	loginData, err := LoginUserInstance(tx, *user, ctx.Value(consts.UserActAsContext).(string))
	if err != nil {
		return nil, err
	}

	history := models.LoginHistory{
		BasicModel: models.BasicModel{
			ID: func() *uuid.UUID {
				id := uuid.New()
				return &id
			}(),
		},
		UserID:    *user.ID,
		TokenID:   &loginData.TokenID,
		UserAgent: &input.UserAgent,
		IP:        &input.IP,
	}

	err = tx.Create(&history).Error
	if err != nil {
		tx.Rollback()
		return nil, errors.New(consts.InternalServerError)
	}

	tx.Commit()

	return loginData, nil
}

func LoginByUsernameOrMobile(ctx context.Context, input appmodels.LoginInputModel) (*appmodels.LoginOutputModel, error) {
	db := dbpkg.MustGormDBConn(ctx)

	user, err := user.GetUserByUsername(input.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(consts.LoginFailed)
		}
		return nil, errors.New(consts.InternalServerError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, errors.New(consts.LoginFailed)
	}
	tx := db.Begin()
	loginData, err := LoginUserInstance(tx, *user, ctx.Value(consts.UserActAsContext).(string))
	if err != nil {
		return nil, err
	}

	history := models.LoginHistory{
		BasicModel: models.BasicModel{
			ID: func() *uuid.UUID {
				id := uuid.New()
				return &id
			}(),
		},
		UserID:    *user.ID,
		TokenID:   &loginData.TokenID,
		UserAgent: &input.UserAgent,
		IP:        &input.IP,
	}

	err = tx.Create(&history).Error
	if err != nil {
		tx.Rollback()
		return nil, errors.New(consts.InternalServerError)
	}

	tx.Commit()

	return loginData, nil
}

func LoginUserInstance(dbConn *gorm.DB, user models.User, actAs string) (*appmodels.LoginOutputModel, error) {

	switch actAs {
	// کاربر میخواهد از روت ادمین لاگین کند
	case consts.UserActAsAdmin:
		{
			if !user.Role.Permitted(models.ACTION_CAN_LOGIN_ADMIN) {
				return nil, errors.New(consts.LoginFailed)
			}
		}
	// کاربر میخواهد از روت کاربر لاگین کند
	case consts.UserActAsUser:
		{
			if !user.Role.Permitted(models.ACTION_CAN_LOGIN_USER) {
				return nil, errors.New(consts.LoginFailed)
			}
		}
	}

	jwtToken := token.NewToken(map[string]interface{}{
		"userID":   user.ID,
		"username": user.Username,
	})

	if !user.Enabled {
		return nil, errors.New(consts.UserIsDisabled)
	}

	tokenStr, err := token.String(jwtToken)

	if err != nil {
		return nil, errors.New(consts.InternalServerError)
	}

	output := &appmodels.LoginOutputModel{
		Token:     tokenStr,
		ExpiresAt: jwtToken.Expiration(),
		ExpiresIn: jwtToken.Expiration().Unix() - time.Now().Unix(),
		User: appmodels.LoginOutputUserModel{
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}

	authToken := models.AuthToken{
		BasicModel: models.BasicModel{
			ID: func() *uuid.UUID {
				id := uuid.MustParse(jwtToken.JwtID())
				return &id
			}(),
		},
		UserID:    *user.ID,
		ExpiresAt: output.ExpiresAt,
		Revoked:   false,
	}
	err = dbConn.Create(&authToken).Error
	if err != nil {
		return nil, errors.New(consts.InternalServerError)
	}

	output.TokenID = *authToken.ID
	return output, nil
}

// RevokeAuthTokenByID
func RevokeAuthTokenByID(db *gorm.DB, id uuid.UUID) error {
	return db.Model(&models.AuthToken{}).
		Where(`"id"`, id.String()).
		UpdateColumn("revoked", true).
		Error
}
