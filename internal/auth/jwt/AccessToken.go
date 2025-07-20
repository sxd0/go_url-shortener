package jwt

import (
	"time"
)

const AccessTokenTTL = time.Minute * 15

func (j *JWT) CreateAccessToken(userID uint, email string) (string, error) {
	return j.createToken(JWTData{
		Email:     email,
		UserID:    userID,
		TokenType: "access",
		Exp:       AccessTokenTTL,
	})
}

func (j *JWT) ParseAccessToken(token string) (bool, *JWTData) {
	ok, data := j.parseToken(token)
	if !ok || data.TokenType != "access" {
		return false, nil
	}
	return true, data
}
