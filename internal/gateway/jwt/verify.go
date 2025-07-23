package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrWrongType    = errors.New("token is not access token")
)

type Verifier struct {
	publicKey []byte
}

func NewVerifier(publicKey string) *Verifier {
	return &Verifier{
		publicKey: []byte(publicKey),
	}
}

func (v *Verifier) ParseToken(tokenStr string) (uint, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(v.publicKey)
	}

	token, err := jwt.Parse(tokenStr, keyFunc, jwt.WithValidMethods([]string{"RS256"}))
	if err != nil || !token.Valid {
		return 0, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidToken
	}

	if expRaw, ok := claims["exp"].(float64); ok {
		if int64(expRaw) < time.Now().Unix() {
			return 0, ErrInvalidToken
		}
	}

	if typ, ok := claims["token_type"].(string); !ok || typ != "access" {
		return 0, ErrWrongType
	}

	idFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, ErrInvalidToken
	}
	return uint(idFloat), nil
}
