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
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		authHeader, ok := md["authorization"]
		if !ok || len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		}

		token := parseBearerToken(authHeader[0])
		if token == "" {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization token format")
		}

		valid, jwtData := jwtService.ParseAccessToken(token)
		if !valid {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired access token")
		}

		ctx = context.WithValue(ctx, ContextJWTKey{}, jwtData)

		return handler(ctx, req)
	}
}

func isPublicMethod(method string) bool {
	return method == "/authpb.AuthService/Login" ||
		method == "/authpb.AuthService/Register" ||
		method == "/authpb.AuthService/Refresh" ||
		method == "/authpb.AuthService/VerifyToken"
}

func parseBearerToken(header string) string {
	if !strings.HasPrefix(header, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(header, "Bearer ")
}
