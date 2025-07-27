package handler

import (
	"github.com/sxd0/go_url-shortener/internal/gateway/jwt"
	"github.com/sxd0/go_url-shortener/proto/gen/go/authpb"
	"github.com/sxd0/go_url-shortener/proto/gen/go/linkpb"
	"github.com/sxd0/go_url-shortener/proto/gen/go/statpb"
)

type Deps struct {
	AuthClient authpb.AuthServiceClient
	LinkClient linkpb.LinkServiceClient
	StatClient statpb.StatServiceClient
	Verifier   *jwt.Verifier
}
