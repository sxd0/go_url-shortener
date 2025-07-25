package gateway

import (
	"github.com/sxd0/go_url-shortener/internal/gateway/authclient"
	"github.com/sxd0/go_url-shortener/internal/gateway/jwt"
	"github.com/sxd0/go_url-shortener/internal/gateway/linkclient"
	"github.com/sxd0/go_url-shortener/proto/gen/go/statpb"
)

type Deps struct {
	AuthClient *authclient.AuthClient
	LinkClient *linkclient.LinkClient
	StatClient statpb.StatServiceClient
	Verifier   *jwt.Verifier
}
