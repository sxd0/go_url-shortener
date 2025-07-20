package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTData struct {
	Email     string
	UserID    uint
	TokenType string
	Exp       time.Duration
}

type JWT struct {
	Secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) createToken(data JWTData) (string, error) {
	claims := jwt.MapClaims{
		"email":      data.Email,
		"user_id":    data.UserID,
		"exp":        time.Now().Add(data.Exp).Unix(),
		"token_type": data.TokenType,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(j.Secret))
}

func (j *JWT) parseToken(token string) (bool, *JWTData) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})

	if err != nil || !t.Valid {
		return false, nil
	}

	claims, ok := t.Claims.(jwt.MapClaims)
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
