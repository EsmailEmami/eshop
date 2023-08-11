package token

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/lestrrat-go/jwx/jwt"
)

func TestLoadTokenFromHttpRequest(t *testing.T) {
	t.Run("no token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "somewhere", nil)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		_, _, err := LoadTokenFromHttpRequest(req)
		if err == nil {
			t.Errorf("LoadTokenFromHttpRequest() wants: %v, got: no error", errors.New("no token found"))
		}
	})
	t.Run("valid token exists in header", func(t *testing.T) {
		os.Setenv("KEYS_PATH", "../../oauth-keys")
		err := InitJWT()
		if err != nil {
			t.Error(err)
		}

		payload := make(map[string]interface{})
		payload[jwt.SubjectKey] = 1
		payload["username"] = "mgh"
		payload["userID"] = 1
		tokenString, _ := NewTokenString(payload)

		req, _ := http.NewRequest(http.MethodGet, "somewhere", nil)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", tokenString)

		_, _, err = LoadTokenFromHttpRequest(req)
		if err != nil {
			t.Errorf("LoadTokenFromHttpRequest() wants: no error, got: %v", err)
		}
	})
	t.Run("valid token exists in query params: Authorization", func(t *testing.T) {
		os.Setenv("KEYS_PATH", "../../oauth-keys")
		err := InitJWT()
		if err != nil {
			t.Error(err)
		}

		payload := make(map[string]interface{})
		payload[jwt.SubjectKey] = 1
		payload["username"] = "mgh"
		payload["userID"] = 1
		tokenString, _ := NewTokenString(payload)

		req, _ := http.NewRequest(http.MethodGet, "somewhere", nil)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		params := req.URL.Query()
		params.Set("Authorization", tokenString)
		req.URL.RawQuery = params.Encode()

		_, _, err = LoadTokenFromHttpRequest(req)
		if err != nil {
			t.Errorf("LoadTokenFromHttpRequest() wants: no error, got: %v", err)
		}
	})
	t.Run("valid token exists in query params: authorization", func(t *testing.T) {
		os.Setenv("KEYS_PATH", "../../oauth-keys")
		err := InitJWT()
		if err != nil {
			t.Error(err)
		}

		payload := make(map[string]interface{})
		payload[jwt.SubjectKey] = 1
		payload["username"] = "mgh"
		payload["userID"] = 1
		tokenString, _ := NewTokenString(payload)

		req, _ := http.NewRequest(http.MethodGet, "somewhere", nil)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		params := req.URL.Query()
		params.Set("authorization", tokenString)
		req.URL.RawQuery = params.Encode()

		_, _, err = LoadTokenFromHttpRequest(req)
		if err != nil {
			t.Errorf("LoadTokenFromHttpRequest() wants: no error, got: %v", err)
		}
	})
	t.Run("valid token exists in cookie: Authorization", func(t *testing.T) {
		os.Setenv("KEYS_PATH", "../../oauth-keys")
		err := InitJWT()
		if err != nil {
			t.Error(err)
		}

		payload := make(map[string]interface{})
		payload[jwt.SubjectKey] = 1
		payload["username"] = "mgh"
		payload["userID"] = 1
		tokenString, _ := NewTokenString(payload)

		req, _ := http.NewRequest(http.MethodGet, "somewhere", nil)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		c := http.Cookie{Name: "Authorization", Value: tokenString, Path: "/"}
		req.AddCookie(&c)

		_, _, err = LoadTokenFromHttpRequest(req)
		if err != nil {
			t.Errorf("LoadTokenFromHttpRequest() wants: no error, got: %v", err)
		}
	})
	t.Run("valid token exists in cookie: authorization", func(t *testing.T) {
		os.Setenv("KEYS_PATH", "../../oauth-keys")
		err := InitJWT()
		if err != nil {
			t.Error(err)
		}

		payload := make(map[string]interface{})
		payload[jwt.SubjectKey] = 1
		payload["username"] = "mgh"
		payload["userID"] = 1
		tokenString, _ := NewTokenString(payload)

		req, _ := http.NewRequest(http.MethodGet, "somewhere", nil)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		c := http.Cookie{Name: "authorization", Value: tokenString, Path: "/"}
		req.AddCookie(&c)

		_, _, err = LoadTokenFromHttpRequest(req)
		if err != nil {
			t.Errorf("LoadTokenFromHttpRequest() wants: no error, got: %v", err)
		}
	})
	t.Run("invalid token", func(t *testing.T) {
		os.Setenv("KEYS_PATH", "../../oauth-keys")
		err := InitJWT()
		if err != nil {
			t.Error(err)
		}

		req, _ := http.NewRequest(http.MethodGet, "somewhere", nil)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		c := http.Cookie{Name: "authorization", Value: "invalid", Path: "/"}
		req.AddCookie(&c)

		_, _, err = LoadTokenFromHttpRequest(req)
		if err == nil {
			t.Errorf("LoadTokenFromHttpRequest() wants: %v, got: no error", errors.New("invalid auth token"))
		}
	})
}
