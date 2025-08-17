package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func newTestJWT(t *testing.T) *JWT {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	return NewJWT(priv, &priv.PublicKey)
}

func TestCreateAndParseAccessToken(t *testing.T) {
	j := newTestJWT(t)
	token, err := j.CreateAccessToken(1, "user@example.com")
	if err != nil {
		t.Fatalf("CreateAccessToken returned error: %v", err)
	}
	ok, data := j.ParseAccessToken(token)
	if !ok {
		t.Fatalf("ParseAccessToken failed")
	}
	if data.Email != "user@example.com" || data.UserID != 1 || data.TokenType != "access" {
		t.Fatalf("unexpected token data: %+v", data)
	}
}

func TestParseAccessTokenWrongType(t *testing.T) {
	j := newTestJWT(t)
	token, err := j.CreateRefreshToken(2, "user@example.com")
	if err != nil {
		t.Fatalf("CreateRefreshToken returned error: %v", err)
	}
	ok, _ := j.ParseAccessToken(token)
	if ok {
		t.Fatalf("expected ParseAccessToken to fail for refresh token")
	}
}

func TestParseTokenInvalid(t *testing.T) {
	j := newTestJWT(t)
	ok, _ := j.ParseAccessToken("invalid.token")
	if ok {
		t.Fatalf("expected parsing invalid token to fail")
	}
}

func TestParseTokenExpired(t *testing.T) {
	j := newTestJWT(t)
	token, err := j.createToken(JWTData{Email: "e", UserID: 1, TokenType: "access", Exp: -1})
	if err != nil {
		t.Fatalf("createToken returned error: %v", err)
	}
	ok, _ := j.ParseAccessToken(token)
	if ok {
		t.Fatalf("expected expired token to be invalid")
	}
}