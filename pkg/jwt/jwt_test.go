package jwt_test

import (
	"go/test-http/pkg/jwt"
	"testing"
)

func TestJWTCreate(t *testing.T) {
	const email = "user@example.com"
	jwtService := jwt.NewJWT("4uCeDM6q0l+2R6JJaQdkJJTGo/f9cPq7MlxJPz463ng=")
	token, err := jwtService.Create(jwt.JWTData{
		Email: email,
	})
	if err != nil {
		t.Fatal(err)
	}
	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatal("Token is invalid")
	}
	if data.Email != email {
		t.Fatalf("Email %s not equal %s", data.Email, email)
	}
}
