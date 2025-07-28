package server

import (
	"context"
	"strings"

	"github.com/sxd0/go_url-shortener/internal/auth/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ContextJWTKey struct{}

func NewJWTUnaryInterceptor(jwtService *jwt.JWT) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if jwtService == nil {
			return nil, status.Error(codes.Internal, "jwt service not initialized")
		}

		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}
		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		}

		token := parseBearer(values[0])
		if token == "" {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization token format")
		}

		valid, data := jwtService.ParseAccessToken(token)
		if !valid {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired access token")
		}

		ctx = context.WithValue(ctx, ContextJWTKey{}, data)
		return handler(ctx, req)
	}
}

func isPublicMethod(full string) bool {
	switch {
	case strings.HasSuffix(full, "/Register"),
		strings.HasSuffix(full, "/Login"),
		strings.HasSuffix(full, "/Refresh"),
		strings.HasSuffix(full, "/VerifyToken"):
		return true
	default:
		return false
	}
}

func parseBearer(h string) string {
	const p = "Bearer "
	if strings.HasPrefix(h, p) {
		return strings.TrimPrefix(h, p)
	}
	return ""
}
