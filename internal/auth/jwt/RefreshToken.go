package jwt

import (
	"time"
)

const RefreshTokenTTL = time.Hour * 24 * 7

func (j *JWT) CreateRefreshToken(userID uint, email string) (string, error) {
	return j.createToken(JWTData{
		Email:     email,
		UserID:    userID,
		TokenType: "refresh",
		Exp:       RefreshTokenTTL,
	})
}

func (j *JWT) ParseRefreshToken(token string) (bool, *JWTData) {
	ok, data := j.parseToken(token)
	if !ok || data.TokenType != "refresh" {
		return false, nil
	}
	return true, data
}
