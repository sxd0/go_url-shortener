package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTData struct {
	Email     string
	UserID    uint
	IsRefresh bool
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

func (j *JWT) Create(data JWTData) (string, error) {
	claims := jwt.MapClaims{
		"email":   data.Email,
		"user_id": data.UserID,
		"exp":     time.Now().Add(data.Exp).Unix(),
	}

	if data.IsRefresh {
		claims["refresh"] = true
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return s, nil
}

func (j *JWT) Parse(token string) (bool, *JWTData) {
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

	data := JWTData{
		Email:  claims["email"].(string),
		UserID: uint(claims["user_id"].(float64)),
	}

	if val, ok := claims["refresh"].(bool); ok && val {
		data.IsRefresh = true
	}

	return true, &data
}
