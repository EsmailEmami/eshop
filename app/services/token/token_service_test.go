package token

import (
	"os"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/jwt"
)

func TestInitJWT(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		os.Setenv("KEYS_PATH", "../../oauth-keys")
		err := InitJWT()
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("no key found", func(t *testing.T) {
		os.Setenv("KEYS_PATH", "invalidpath")
		err := InitJWT()
		if err == nil {
			t.Error("error wanted but got nil")
		}
	})
}

func TestNewToken(t *testing.T) {
	os.Setenv("KEYS_PATH", "../../oauth-keys")
	err := InitJWT()
	if err != nil {
		t.Error(err)
	}
	payload := make(map[string]interface{})
	payload[jwt.SubjectKey] = "mgh"
	token := NewToken(payload)
	s, ok := token.Get(jwt.SubjectKey)
	if !ok {
		t.Error("unable to fetch subject")
	}

	if s != "mgh" {
		t.Fatal("invalid subject")
	}
}

func TestString(t *testing.T) {
	os.Setenv("KEYS_PATH", "../../oauth-keys")
	err := InitJWT()
	if err != nil {
		t.Error(err)
	}
	payload := make(map[string]interface{})
	payload[jwt.SubjectKey] = "mgh"
	token := NewToken(payload)
	_, err = String(token)
	if err != nil {
		t.Errorf("String() wants error: nil, got: %v", err)
	}
}

func TestNewTokenString(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		os.Setenv("KEYS_PATH", "../../oauth-keys")
		err := InitJWT()
		if err != nil {
			t.Error(err)
		}
		payload := make(map[string]interface{})
		payload[jwt.SubjectKey] = "mgh"
		token, err := NewTokenString(payload)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(token)
	})
}

func TestParseToken(t *testing.T) {
	os.Setenv("KEYS_PATH", "../../oauth-keys")
	err := InitJWT()
	if err != nil {
		t.Error(err)
	}

	t.Run("invalid token", func(t *testing.T) {
		_, err := ParseToken("invalid", true)
		if err == nil {
			t.Fatal("error wanted but got nil")
		}
	})

	t.Run("expired token", func(t *testing.T) {
		payload := make(map[string]interface{})
		payload[jwt.ExpirationKey] = time.Now().Add(-100 * time.Minute).UTC().Unix()

		token, _ := NewTokenString(payload)
		_, err := ParseToken(token, true)
		if err == nil {
			t.Fatal("error wanted but got nil")
		}
	})

	t.Run("ok", func(t *testing.T) {
		payload := make(map[string]interface{})

		tokenString, _ := NewTokenString(payload)
		t.Log(tokenString)
		_, err := ParseToken(tokenString, true)
		if err != nil {
			t.Fatal("error on parsing token:", err)
		}
	})
}

func BenchmarkParseToken(b *testing.B) {
	os.Setenv("KEYS_PATH", "../../oauth-keys")
	_ = InitJWT()
	payload := make(map[string]interface{})
	payload[jwt.SubjectKey] = "MGH"
	tokenString, _ := NewTokenString(payload)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ParseToken(tokenString, true)
	}
}
