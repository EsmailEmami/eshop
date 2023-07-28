package token

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"gorm.io/gorm"
)

func LoadTokenFromHttpRequest(req *http.Request) (token *jwt.Token, tokenString string, err error) {
	c, err := req.Cookie("Authorization")
	if err == nil {
		tokenString = c.Value
	}

	if tokenString == "" {
		tokenString = req.Header.Get("Authorization")
	}

	if tokenString == "" {
		tokenString = req.URL.Query().Get("Authorization")
		if tokenString == "" {
			tokenString = req.URL.Query().Get("authorization")
		}
	}

	if tokenString == "" {
		c, err := req.Cookie("authorization")
		if err == nil {
			tokenString = c.Value
		}
	}

	if tokenString == "" {
		return nil, "", errors.New("no token found")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err = ParseToken(tokenString, true)
	if err != nil {
		return nil, "", errors.New("invalid auth token")
	}

	return token, tokenString, nil
}

func GetValidAuthTokenByID(db *gorm.DB, tokenID uuid.UUID) (*models.AuthToken, *models.User, error) {
	var authToken models.AuthToken

	if err := db.Where(`"id"=?`, tokenID).Where(`"revoked"=?`, false).
		Where(`"expires_at">?`, time.Now()).First(&authToken).Error; err != nil {
		return nil, nil, err
	}

	var user models.User

	if err := db.Where(`"id"=?`, authToken.UserID).Preload("Role").First(&user).Error; err != nil {
		return nil, nil, err
	}

	return &authToken, &user, nil
}
