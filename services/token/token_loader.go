package token

import (
	"errors"
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/jwt"
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
