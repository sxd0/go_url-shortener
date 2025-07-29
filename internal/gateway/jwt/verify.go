package jwt

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrWrongType    = errors.New("token is not access token")
	ErrExpired      = errors.New("token expired")
	ErrWrongMethod  = errors.New("invalid signing method")
)

type AccessClaims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type Verifier struct {
	publicKey *rsa.PublicKey
}

func NewVerifier(pem string) *Verifier {
	pk, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
	if err != nil {
		panic("jwt: bad RSA public key: " + err.Error())
	}
	return &Verifier{publicKey: pk}
}

func (v *Verifier) ParseToken(tokenStr string) (uint, error) {
	keyFunc := func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, ErrWrongMethod
		}
		return v.publicKey, nil
	}

	var claims AccessClaims
	tok, err := jwt.ParseWithClaims(
		tokenStr,
		&claims,
		keyFunc,
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
	)
	if err != nil || !tok.Valid {
		return 0, ErrInvalidToken
	}
	if claims.TokenType != "access" {
		return 0, ErrWrongType
	}
	if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		return 0, ErrExpired
	}
	return uint(claims.UserID), nil
}
