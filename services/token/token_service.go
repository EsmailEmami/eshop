package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// and load the signKey on every signing request? depends on  your usage i guess
var (
	privateKey *rsa.PrivateKey
)

func InitJWT() (err error) {
	keyPath := viper.GetString("keys.private")
	bts, err := os.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(bts)

	pk, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	prvKey, ok := pk.(*rsa.PrivateKey)
	if !ok {
		return errors.New("Invalid private key")
	}

	privateKey = prvKey

	return nil
}

func NewToken(payload map[string]interface{}) jwt.Token {
	token := jwt.New()

	// IssuedAtKey and NotBeforeKey and ExpirationKey could be override by payload
	_ = token.Set(jwt.IssuedAtKey, time.Now().Add(time.Second*-1).UTC())
	_ = token.Set(jwt.NotBeforeKey, time.Now().Add(time.Second*-1).UTC())
	tokenExp := time.Now().UTC().Add(24 * 29 * time.Hour)
	_ = token.Set(jwt.ExpirationKey, tokenExp)

	_ = token.Set(jwt.JwtIDKey, uuid.New().String())

	for k, v := range payload {
		_ = token.Set(k, v)
	}

	return token
}

func String(token jwt.Token) (tokenString string, err error) {
	bts, err := jwt.Sign(token, jwa.RS256, privateKey)
	if err != nil {
		logrus.Errorf("failed to generate signed payload: %s\n", err)
		return
	}
	return string(bts), nil
}

func NewTokenString(payload map[string]interface{}) (tokenString string, err error) {
	token := NewToken(payload)
	bts, err := jwt.Sign(token, jwa.RS256, privateKey)
	if err != nil {
		logrus.Errorf("failed to generate signed payload: %s\n", err)
		return
	}
	return string(bts), nil
}

func ParseToken(tokenString string, validateToken bool) (*jwt.Token, error) {
	token, err := jwt.ParseString(
		tokenString,
		jwt.WithVerify(jwa.RS256, &privateKey.PublicKey),
		jwt.WithValidate(validateToken),
	)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
