package jwt

import (
	"crypto/rsa"
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type JWTData struct {
	Email     string
	UserID    uint
	TokenType string
	Exp       time.Duration
}

type JWT struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewJWT(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *JWT {
	return &JWT{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}

func (j *JWT) createToken(data JWTData) (string, error) {
	claims := jwtlib.MapClaims{
		"email":      data.Email,
		"user_id":    data.UserID,
		"exp":        time.Now().Add(data.Exp).Unix(),
		"token_type": data.TokenType,
	}

	t := jwtlib.NewWithClaims(jwtlib.SigningMethodRS256, claims)
	return t.SignedString(j.PrivateKey)
}

func (j *JWT) parseToken(token string) (bool, *JWTData) {
	t, err := jwtlib.Parse(token, func(t *jwtlib.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtlib.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.PublicKey, nil
	})

	if err != nil || !t.Valid {
		return false, nil
	}

	claims, ok := t.Claims.(jwtlib.MapClaims)
	if !ok {
		return false, nil
	}

	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) < time.Now().Unix() {
		return false, nil
	}

	tokenType, ok := claims["token_type"].(string)
	if !ok {
		return false, nil
	}

	return true, &JWTData{
		Email:     claims["email"].(string),
		UserID:    uint(claims["user_id"].(float64)),
		TokenType: tokenType,
	}
}
